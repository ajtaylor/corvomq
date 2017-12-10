package main

import (
	"flag"
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	var port = flag.Int("port", 8000, "Port to listen on")

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)

	// Routes
	e.File("/*", "html/app.html")

	e.Static("/static", "static")

	// Start server
	log.Println("Starting server on localhost:8080...")
	e.Run(standard.New(":8080"))
}
