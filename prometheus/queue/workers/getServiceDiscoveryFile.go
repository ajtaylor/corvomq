package workers

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	minio "github.com/minio/minio-go"
	"github.com/nats-io/nats"
)

var (
	spacesAccessKey   = os.Getenv("SPACES_ACCESS_KEY")
	spacesSecurityKey = os.Getenv("SPACES_SECURITY_KEY")
)

// GetServiceDiscoveryFileMessage - processed by this worker
type GetServiceDiscoveryFileMessage struct {
	ExporterType string
	Filename     string
}

// RunGetServiceDiscoveryFileSubscriber creates subscriber
func RunGetServiceDiscoveryFileSubscriber() {
	var sdFileMsg GetServiceDiscoveryFileMessage

	logger := log.WithFields(log.Fields{
		"worker":   "GetServiceDiscoveryFile",
		"function": "RunGetServiceDiscoveryFileSubscriber",
	})

	opts := nats.DefaultOptions
	opts.Url = "nats://88.202.188.57:4222"
	nc, err := opts.Connect()
	if err != nil {
		logger.WithFields(log.Fields{
			"error":       err.Error(),
			"nats_server": opts.Url,
		}).Fatal("Connect to NATS server")
	}
	logger.WithField("nats_server", opts.Url).Info("Connected to NATS server")

	nc.QueueSubscribe("NewSDFile", "SDFile_Queue", func(msg *nats.Msg) {
		logger.WithField("msg", msg.Data).Info("Received msg")

		err = json.Unmarshal(msg.Data, &sdFileMsg)
		if err != nil {
			logger.WithFields(log.Fields{
				"msg":   &sdFileMsg,
				"error": err.Error(),
			}).Error("Unmarshal msg")
		}

		// Connect to Spaces
		spacesClient, err := minio.New("nyc3.digitaloceanspaces.com", spacesAccessKey, spacesSecurityKey, true)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Connect to Spaces")
		}

		// Get service discovery file
		err = spacesClient.FGetObject("corvomq", "service_discovery/"+sdFileMsg.Filename, "/opt/prometheus-1.4.1.linux-amd64/service_discovery/"+sdFileMsg.ExporterType+"_exporter/"+sdFileMsg.Filename, minio.GetObjectOptions{})
		if err != nil {
			logger.WithFields(log.Fields{
				"filename": sdFileMsg.Filename,
				"error":    err.Error(),
			}).Error("Get service discovery file from Spaces")
		}
	})
}
