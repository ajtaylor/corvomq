package resources

import (
	"net/http"

	"github.com/ajtaylor/corvomq/api/repository"
)

// Five waits 5 seconds
func Five(w http.ResponseWriter, r *http.Request) {
	var n int
	_ = repository.Connection.QueryRow("SELECT pg_sleep(5);").Scan(&n)
	w.Write([]byte("5 secs"))
}

// Ten waits 10 seconds
func Ten(w http.ResponseWriter, r *http.Request) {
	var n int
	_ = repository.Connection.QueryRow("SELECT pg_sleep(10);").Scan(&n)
	w.Write([]byte("10 secs"))
}
