package repository

import (
	"log"
	"os"
	"time"

	"github.com/jackc/pgx"
)

// Connection is a pgx ConnPool
var Connection *pgx.ConnPool

// Configure sets up the Connection pgx ConnPool
func Configure() {
	var (
		pgxConfig pgx.ConnConfig
		pgxPool   *pgx.ConnPool
		err       error
		n         time.Time
	)

	pgxConfig, err = pgx.ParseURI(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Printf("Can't parse DB URI: %v", err.Error())
	}

	pgxPool, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgxConfig})
	if err != nil {
		log.Printf("Can't create DB Connection pool: %v", err.Error())
	}

	err = pgxPool.QueryRow("SELECT Now() as Now;").Scan(&n)
	if err != nil {
		log.Printf("Can't execute DB check: %v", err.Error())
	}

	log.Printf("DB initialised at %v", n)

	Connection = pgxPool
}
