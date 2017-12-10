package resources

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ajtaylor/corvomq/api/repository"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
)

// User is a CorvoMQ user
type User struct {
	ID             int32  `json:"id"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	EmailAddress   string `json:"emailAddress"`
	Password       string `json:"password"`
	BroadcastKey   string `json:"broadcast_key"`
	HashedPassword []byte `json:"hashed_password"`
}

// LoginRequest is
type LoginRequest struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

// Bind binds the JSON from the LoginRequest
func (l *LoginRequest) Bind(r *http.Request) error {
	return nil
}

// Login processes a user login
func Login(w http.ResponseWriter, r *http.Request) {
	var appuserID int32
	var hashedPassword []byte
	var tokenAuth *jwtauth.JwtAuth

	logger := lg.WithFields(logrus.Fields{
		"module":   "resources.user",
		"function": "Login",
	})

	loginData := &LoginRequest{}

	if err := render.Bind(r, loginData); err != nil {
		logger.WithField("error", err.Error()).Error("Bind request data")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err := repository.Connection.QueryRow("SELECT appuser_id, hashed_password FROM corvomq.login_appuser($1);",
		loginData.EmailAddress).Scan(&appuserID, &hashedPassword)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Execute corvomq.login_appuser")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(loginData.Password))
	if err != nil {
		logger.WithField("error", err.Error()).Error("Password compare")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	tokenAuth = jwtauth.New("HS256", []byte("0nTheBu5e5"), nil)
	jwtauth.ExpireIn(time.Second * 3600 * 1)
	_, tokenString, _ := tokenAuth.Encode(jwtauth.Claims{"user_id": appuserID, "iat": jwtauth.EpochNow()})
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"token": tokenString})
}
