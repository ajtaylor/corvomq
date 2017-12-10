package resources

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/ajtaylor/corvomq/api/repository"
	"github.com/go-chi/render"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

// RegistrationRequest defines a new registration
type RegistrationRequest struct {
	Firstname        string `form:"firstname" json:"firstname"`
	Lastname         string `form:"lastname" json:"lastname"`
	EmailAddress     string `form:"email_address" json:"email_address"`
	Password         string `json:"password"`
	OrganisationName string `json:"organisation_name"`
	HashedPassword   []byte
}

// Bind binds
func (reg *RegistrationRequest) Bind(r *http.Request) error {
	return nil
}

// CreateRegistration saves a new registration to the database
func CreateRegistration(w http.ResponseWriter, r *http.Request) {
	var registration = &RegistrationRequest{}

	logger := lg.WithFields(logrus.Fields{
		"module":   "resources.registration",
		"function": "CreateRegistration",
	})

	err := render.Bind(r, registration)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Bind request data")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var appuserID int32
	var organisationID int32
	var broadcastKey = uuid4()

	err = repository.Connection.QueryRow("SELECT appuser_id, organisation_id FROM corvomq.register_organisation($1, $2, $3, $4, $5, $6);",
		registration.Firstname,
		registration.Lastname,
		registration.EmailAddress,
		registration.HashedPassword,
		registration.OrganisationName,
		broadcastKey).Scan(&appuserID, &organisationID)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Execute corvomq.register_organisation")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]bool{"Created": true})
}

func uuid4() string {
	logger := lg.WithFields(logrus.Fields{
		"module":   "resources.registration",
		"function": "uuid4",
	})

	b := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Create random bytes")
	}
	if n != len(b) || err != nil {
		return ""
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
