package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/nats-io/nats"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	shortid "github.com/ventu-io/go-shortid"

	queue "github.com/ajtaylor/corvomq/queue/workers"

	"github.com/ajtaylor/corvomq/api/repository"
	"github.com/go-chi/render"
)

// Environment describes a UserMQ environment
type Environment struct {
	ID             int       `json:"id"`
	OrganisationID int       `json:"organisation_id"`
	Name           string    `json:"name"`
	Server         string    `json:"server"`
	Infrastructure string    `json:"infrastructure"`
	URL            string    `json:"url"`
	TLSEnabled     bool      `json:"tls_enabled"`
	CreatedAt      time.Time `json:"created_at"`
}

// Bind binds the JSON from the EnvironmentRequest
func (e *Environment) Bind(r *http.Request) error {
	return nil
}

// GetEnvironmentList returns a JSON list
func GetEnvironmentList(w http.ResponseWriter, r *http.Request) {
	var id int
	var name string
	var server string
	var infrastructure string
	var url string
	var tlsEnabled bool
	var createdAt time.Time
	var environmentList []Environment

	logger := lg.WithFields(logrus.Fields{
		"module":   "resources.environment",
		"function": "GetEnvironmentList",
	})

	_, claims, _ := jwtauth.FromContext(r.Context())
	appuserID := int(claims["user_id"].(float64))

	rows, err := repository.Connection.Query("SELECT id, name, server, infrastructure, url, tls_enabled, created_at FROM corvomq.appuser_get_environments($1);", appuserID)
	defer rows.Close()
	if err != nil {
		logger.WithField("error", err.Error()).Error("Execute corvomq.appuser_get_environments")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	for rows.Next() {
		err := rows.Scan(&id, &name, &server, &infrastructure, &url, &tlsEnabled, &createdAt)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Read row from corvomq.appuser_get_environments")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}
		environmentList = append(environmentList, Environment{ID: id,
			Name:           name,
			Server:         server,
			Infrastructure: infrastructure,
			URL:            url,
			TLSEnabled:     tlsEnabled,
			CreatedAt:      createdAt})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string][]Environment{"environments": environmentList})
}

// CreateEnvironment saves new user environment to DB and sends create message to MQ server
func CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var environment = &Environment{}
	var environmentID int

	logger := lg.WithFields(logrus.Fields{
		"module":   "resources.environment",
		"function": "CreateEnvironment",
	})

	_, claims, _ := jwtauth.FromContext(r.Context())
	appuserID := int(claims["user_id"].(float64))

	// User name must start with with an alpha character
	startsWithAlpha, _ := regexp.Compile("^[a-zA-Z]")

	// Build URL for new user environment
	urlShortid, _ := shortid.Generate()
	// Ensure URL starts with an alpha character
	for !startsWithAlpha.MatchString(urlShortid) {
		urlShortid, _ = shortid.Generate()
	}
	urlInt := fmt.Sprintf("%d", time.Now().Nanosecond())
	urlsegment := urlShortid + "-" + urlInt

	if err := render.Bind(r, environment); err != nil {
		logger.WithField("error", err.Error()).Error("Bind request data")
		render.Render(w, r, ErrInvalidRequest(err))
	}

	// Save to DB
	err := repository.Connection.QueryRow("SELECT environment_id FROM corvomq.create_environment($1, $2, $3, $4, $5, $6);",
		appuserID,
		environment.Name,
		environment.Server,
		environment.Infrastructure,
		urlsegment,
		environment.TLSEnabled).Scan(&environmentID)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Execute corvomq.create_environment")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	vmMsg := queue.CreateVirtualMachineMsg{
		EnvironmentID:  environmentID + appuserID,
		Infrastructure: environment.Infrastructure,
		Server:         environment.Server,
		TLSEnabled:     environment.TLSEnabled,
		Host:           urlsegment,
	}

	natsOptions := nats.DefaultOptions
	natsOptions.Url = "nats://88.202.188.57:4222"
	nc, err := natsOptions.Connect()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":       err.Error(),
			"nats_server": natsOptions.Url,
		}).Error("Connect to STAN server")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	defer nc.Close()

	msg, _ := json.Marshal(vmMsg)

	log.Println("Sending to STAN server...")
	err = nc.Publish("CreateVM", msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":       err.Error(),
			"nats_server": natsOptions.Url,
			"message":     msg,
		}).Error("Send message to STAN server")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]int{"environmentId": environmentID})
}
