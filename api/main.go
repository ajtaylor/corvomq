package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/ajtaylor/corvomq/api/repository"
	"github.com/ajtaylor/corvomq/api/resources"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pressly/lg"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {
	var tokenAuth *jwtauth.JwtAuth
	var port = flag.Int("port", 8080, "Port to listen on")

	flag.Parse()

	repository.Configure()

	tokenAuth = jwtauth.New("HS256", []byte("0nTheBu5e5"), nil)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://app.corvomq.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	})

	infoWriter, err := rotatelogs.New("/var/log/corvomq/api/info.%Y%m%d",
		rotatelogs.WithClock(rotatelogs.UTC),
		rotatelogs.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Println(err.Error())
	}

	errorWriter, err := rotatelogs.New("/var/log/corvomq/api/error.%Y%m%d",
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
	r.Use(cors.Handler)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("tta"))
	})

	r.Post("/organisation/register", resources.CreateRegistration)
	r.Post("/user/login", resources.Login)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/environments", resources.GetEnvironmentList)
		r.Post("/environments/create", resources.CreateEnvironment)
	})

	lg.WithField("port", *port).Info("API server started")

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs)

	go func() {
		s := <-sigs
		lg.WithField("signal", s).Info("Received signal")
		appExit()
		os.Exit(1)
	}()

	http.ListenAndServe(":"+strconv.Itoa(*port), r)
}

func appExit() {
	lg.Info("Exiting")
}
