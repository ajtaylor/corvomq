package workers

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
	"regexp"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"

	prometheusQueue "github.com/ajtaylor/corvomq/prometheus/queue/workers"

	minio "github.com/minio/minio-go"
	"github.com/nats-io/nats"
	"github.com/parnurzeal/gorequest"
	shortid "github.com/ventu-io/go-shortid"
)

const (
	passwordLength          = 15
	messageServerTemplateID = 8921
	ssdVpsPlan              = "0"
	cloudID                 = 120 // Amsterdam
	domain                  = "corvomq.com"
	domainID                = 171664
	apiCreateVMURL          = "https://api.vps.net/ssd_virtual_machines.api10json"
	apiCreateDNSURL         = "https://api.vps.net/domains/171664/records.api10json"
	apiUsername             = "antony.taylor@gmail.com"
	natsServerURL           = "nats://88.202.188.57:4222"
)

var (
	spacesAccessKey   = os.Getenv("SPACES_ACCESS_KEY")
	spacesSecurityKey = os.Getenv("SPACES_SECURITY_KEY")
	apiToken          = os.Getenv("VPS_NET_API_TOKEN")
)

// CreateVirtualMachineMsg - processed by this worker
type CreateVirtualMachineMsg struct {
	TLSEnabled     bool   `json:"tls_enabled"`
	EnvironmentID  int    `json:"environment_id"`
	Server         string `json:"server"`
	Infrastructure string `json:"infrastructure"`
	Host           string `json:"host"`
}

type natsConfig struct {
	Username   string
	Password   string
	TLSEnabled bool
}

type virtualMachineDefinition struct {
	SsdVpsPlan          string           `json:"ssd_vps_plan"`
	Fqdn                string           `json:"fqdn"`
	Tag                 string           `json:"tag"`
	SystemTemplateID    int32            `json:"system_template_id"`
	CloudID             int32            `json:"cloud_id"`
	BackupsEnabled      bool             `json:"backups_enabled"`
	RsyncBackupsEnabled bool             `json:"rsync_backups_enabled"`
	Licences            map[string]int32 `json:"licences"`
}

type createVirtualMachineRequest struct {
	VirtualMachine virtualMachineDefinition `json:"virtual_machine"`
}

type domainRecordDefinition struct {
	TTL  time.Duration `json:"ttl"`
	Data string        `json:"data"`
	Type string        `json:"type"`
	Host string        `json:"host"`
}

type createDomainRecordRequest struct {
	DomainRecord domainRecordDefinition `json:"domain_record"`
}

// createSsdVMResponse is the response from the API Create VM call
type createSsdVMResponse struct {
	VirtualMachine struct {
		Hostname string `json:"hostname"`
		// ID               int    `json:"id"`
		// SSDVPSLabel      string `json:"ssd_vps_label"`
		PrimaryIPAddress struct {
			IPAddress struct {
				IPAddress string `json:"ip_address"`
			} `json:"ip_address"`
		} `json:"primary_ip_address"`
	} `json:"virtual_machine"`
}

// serviceDiscoveryConfig is for the Prometheus server
type serviceDiscoveryConfig struct {
	Hostname       string
	Port           string
	OrganisationID int
}

// RunCreateMQVMSubscriber runs CreateMQVM subscriber
func RunCreateMQVMSubscriber() {
	logger := log.WithFields(log.Fields{
		"worker":   "CreateMQVM",
		"function": "RunCreateMQVMSubscriber",
	})

	logger.Info("Starting CreateMQVM subscriber")

	var createVMMsg CreateVirtualMachineMsg
	var createVM createVirtualMachineRequest
	var vmDef virtualMachineDefinition

	opts := nats.DefaultOptions
	opts.Url = natsServerURL
	nc, err := opts.Connect()
	if err != nil {
		logger.WithFields(log.Fields{
			"error":       err.Error(),
			"nats_server": opts.Url,
		}).Fatal("Connect to NATS server")
	}
	logger.Info("Connected to NATS server")

	// Initialise API request to create new VM
	request := gorequest.New().
		SetBasicAuth(apiUsername, apiToken).
		Timeout(time.Second * 60)

	nc.QueueSubscribe("CreateVM", "CreateVM_Queue", func(msg *nats.Msg) {
		logger.WithField("msg", msg.Data).Info("Received msg")

		err = json.Unmarshal(msg.Data, &createVMMsg)
		if err != nil {
			logger.WithFields(log.Fields{
				"msg":   &createVMMsg,
				"error": err.Error(),
			}).Error("Unmarshal msg")
		}

		// User name must start with with an alpha character
		startsWithAlpha, _ := regexp.Compile("^[a-zA-Z]")

		// Initialise shortid
		sid, _ := shortid.New(1, shortid.DefaultABC, 2343)
		shortid.SetDefault(sid)

		host := createVMMsg.Host
		fqdn := host + "." + domain

		// Generate NATS username and password
		username, _ := shortid.Generate()
		// Ensure username starts with an alpha character
		for !startsWithAlpha.MatchString(username) {
			username, _ = shortid.Generate()
		}
		password := generatePassword()
		passwordEncrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithField("error", err.Error()).Fatalf("Produce bcrypt hash")
			return
		}

		logger.WithFields(log.Fields{
			"username": username,
			"password": password,
		}).Info("NATS auth")

		// Build cloud-init user-data file
		filename := fqdn + "_user-data"

		natsConfig := natsConfig{username, string(passwordEncrypted), createVMMsg.TLSEnabled}

		f, err := os.Create(filename)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filename,
				"error":    err.Error(),
			}).Error("Error creating file")
		}
		defer f.Close()

		// Initialise and execute user-data template
		t, _ := template.ParseFiles("../usermq/cloud-init/user-data/usermq")

		err = t.Execute(f, natsConfig)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": f.Name,
				"error":    err.Error(),
			}).Error("Error applying template")
		}

		// Connect to Spaces
		spacesClient, err := minio.New("nyc3.digitaloceanspaces.com", spacesAccessKey, spacesSecurityKey, true)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Connect to Spaces")
		}

		// Write user-data to Spaces
		_, err = spacesClient.FPutObject("corvomq", "cloud-init/"+filename, filename, minio.PutObjectOptions{})
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filename,
				"error":    err.Error(),
			}).Error("Error writing to Spaces")
		}

		// Build cloud-init meta-data file
		filename = fqdn + "_meta-data"

		f, err = os.Create(filename)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filename,
				"error":    err.Error(),
			}).Error("Error creating file")
		}

		// Initialise and execute meta-data template
		t, _ = template.ParseFiles("../usermq/cloud-init/meta-data")

		err = t.Execute(f, createVMMsg)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": f.Name,
				"error":    err.Error(),
			}).Error("Error applying template")
		}

		// Write meta-data to Spaces
		_, err = spacesClient.FPutObject("corvomq", "cloud-init/"+filename, filename, minio.PutObjectOptions{})
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filename,
				"error":    err.Error(),
			}).Error("Error writing to Spaces")
		}

		// Build Prometheus service discovery files
		// Initialise service discovery template
		t, _ = template.ParseFiles("../prometheus/service_discovery/service_discovery-template.tmpl")

		// Build node service discovery
		filenameNode := "sd_node_" + fqdn + ".json"

		sdConfig := serviceDiscoveryConfig{
			Hostname:       fqdn,
			Port:           "9100",
			OrganisationID: 1}

		fNode, err := os.Create(filenameNode)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filenameNode,
				"error":    err.Error(),
			}).Error("Error creating file")
		}

		err = t.Execute(fNode, sdConfig)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": fNode.Name,
				"error":    err.Error(),
			}).Error("Error applying template")
		}
		fNode.Close()

		// Write node service discovery to Spaces
		_, err = spacesClient.FPutObject("corvomq", "service_discovery/"+filenameNode, filenameNode, minio.PutObjectOptions{})
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filenameNode,
				"error":    err.Error(),
			}).Error("Error writing to Spaces")
		}

		// Build nats service discovery
		filenameNats := "sd_nats_" + fqdn + ".json"

		sdConfig.Port = "7777"

		fNats, err := os.Create(filenameNats)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filenameNats,
				"error":    err.Error(),
			}).Error("Error creating file")
		}

		err = t.Execute(fNats, sdConfig)
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": fNats.Name,
				"error":    err.Error(),
			}).Error("Error applying template")
		}
		fNats.Close()

		// Write nats service discovery to Spaces
		_, err = spacesClient.FPutObject("corvomq", "service_discovery/"+filenameNats, filenameNats, minio.PutObjectOptions{})
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": filenameNats,
				"error":    err.Error(),
			}).Error("Error writing to Spaces")
		}

		fmsg := prometheusQueue.GetServiceDiscoveryFileMessage{ExporterType: "node", Filename: filenameNode}
		prometheusMsg, _ := json.Marshal(fmsg)

		err = nc.Publish("NewSDFile", prometheusMsg)
		if err != nil {
			logger.WithFields(log.Fields{
				"message":    msg,
				"natsServer": natsServerURL,
				"filename":   filenameNode,
				"error":      err.Error(),
			}).Error("Error sending message to NATS server")
		}

		fmsg = prometheusQueue.GetServiceDiscoveryFileMessage{ExporterType: "nats", Filename: filenameNats}
		prometheusMsg, _ = json.Marshal(fmsg)

		err = nc.Publish("NewSDFile", prometheusMsg)
		if err != nil {
			logger.WithFields(log.Fields{
				"message":    msg,
				"natsServer": natsServerURL,
				"filename":   filenameNode,
				"error":      err.Error(),
			}).Error("Error sending message to NATS server")
		}

		// Create VM create API request
		vmDef = virtualMachineDefinition{
			SsdVpsPlan:          ssdVpsPlan,
			Fqdn:                fqdn,
			SystemTemplateID:    messageServerTemplateID,
			CloudID:             cloudID,
			BackupsEnabled:      false,
			RsyncBackupsEnabled: false,
			Licences:            make(map[string]int32),
		}

		createVM = createVirtualMachineRequest{VirtualMachine: vmDef}

		logger.WithFields(log.Fields{
			"createVMRequest": createVM,
		}).Info("CreateVM")

		// Send API request and call handleCreateResponse on completion
		request.Post(apiCreateVMURL).
			Send(createVM).
			End(handleCreateVMResponse)
	})
}

func handleCreateVMResponse(res gorequest.Response, body string, errs []error) {
	logHandleCreateVMResponse := log.WithFields(log.Fields{
		"worker":   "CreateMQVM",
		"function": "handleCreateVMResponse",
	})

	// Check for errors in Create response
	if errs != nil {
		logHandleCreateVMResponse.WithFields(log.Fields{
			"statusCode": res.StatusCode,
			"body":       body,
		}).Error("Checking response")
		for i := 0; i < len(errs); i++ {
			logHandleCreateVMResponse.WithField("error", errs[i].Error()).Error("Response error")
		}
	}

	// Deconstruct response body to get VM info
	bodyBlob := []byte(body)
	var ssdVM createSsdVMResponse
	err := json.Unmarshal(bodyBlob, &ssdVM)
	if err != nil {
		logHandleCreateVMResponse.WithField("error", err.Error()).Error("JSON unmarshal")
		return
	}
	vmIPAddress := ssdVM.VirtualMachine.PrimaryIPAddress.IPAddress.IPAddress
	vmHostname := ssdVM.VirtualMachine.Hostname

	// Initialise API request to create new Domain Record
	request := gorequest.New().
		SetBasicAuth(apiUsername, apiToken).
		Timeout(time.Second * 60)

	domainRecordDef := domainRecordDefinition{
		TTL:  15 * time.Minute / time.Second,
		Data: vmIPAddress,
		Type: "a",
		Host: vmHostname,
	}

	createDomainRecord := createDomainRecordRequest{DomainRecord: domainRecordDef}

	request.Post(apiCreateDNSURL).
		Send(createDomainRecord).
		End(handleCreateDNSResponse)
}

func handleCreateDNSResponse(res gorequest.Response, body string, errs []error) {
	logHandleCreateDNSResponse := log.WithFields(log.Fields{
		"worker":   "CreateMQVM",
		"function": "handleCreateDNSResponse",
	})

	// Check for errors in Create response
	if errs != nil {
		logHandleCreateDNSResponse.WithFields(log.Fields{
			"statusCode": res.StatusCode,
			"body":       body,
		}).Error("Checking response")
		for i := 0; i < len(errs); i++ {
			logHandleCreateDNSResponse.WithField("error", errs[i].Error()).Error("Response error")
		}
	}
}

func generatePassword() string {
	logGeneratePassword := log.WithFields(log.Fields{
		"worker":   "CreateMQVM",
		"function": "generatePassword",
	})

	var ch = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$#%^&*()")
	b := make([]byte, passwordLength)
	max := big.NewInt(int64(len(ch)))
	for i := range b {
		ri, err := rand.Int(rand.Reader, max)
		if err != nil {
			logGeneratePassword.WithField("error", err.Error()).Fatal("Generate random integer")
		}
		b[i] = ch[int(ri.Int64())]
	}
	return string(b)
}
