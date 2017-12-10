package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ajtaylor/corvomq/api/resources"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pressly/lg"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {
	var port = flag.Int("port", 8000, "Port to listen on")

	flag.Parse()

	infoWriter, err := rotatelogs.New("/var/log/corvomq/www/info.%Y%m%d",
		rotatelogs.WithClock(rotatelogs.UTC),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Println(err.Error())
	}

	errorWriter, err := rotatelogs.New("/var/log/corvomq/www/error.%Y%m%d",
		rotatelogs.WithClock(rotatelogs.UTC),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Println(err.Error())
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Hooks.Add(lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  infoWriter,
		logrus.ErrorLevel: errorWriter,
	}))

	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger

	r := chi.NewRouter()

	// r.Use(middleware.DefaultLogger)
	r.Use(lg.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Post("/register-account", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		password := []byte(r.FormValue("password"))
		hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		reg := &resources.RegistrationRequest{}
		reg.Firstname = r.FormValue("firstname")
		reg.Lastname = r.FormValue("lastname")
		reg.OrganisationName = r.FormValue("organisation_name")
		reg.EmailAddress = r.FormValue("email_address")
		reg.HashedPassword = hashedPassword
		b, _ := json.Marshal(reg)
		http.Post("http://localhost:8081/organisation/register", "application/json", bytes.NewBuffer(b))
		http.Redirect(w, r, "http://localhost/thankyou", http.StatusSeeOther)
	})

	lg.WithField("port", *port).Info("WWW server started")

	http.ListenAndServe(":"+strconv.Itoa(*port), r)
}
