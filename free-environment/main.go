package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/jackc/pgx"
	"github.com/ventu-io/go-shortid"
)

func main() {
	sid, _ := shortid.New(1, shortid.DefaultABC, 2343)
	shortid.SetDefault(sid)

	pgxcfg, err := pgx.ParseURI(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Println(err.Error())
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: pgxcfg})
	if err != nil {
		log.Println(err.Error())
	}

	_, _ = pgxpool.Exec("TRUNCATE TABLE corvomq.environment_free_config;")
	_, _ = pgxpool.Exec("ALTER SEQUENCE corvomq.environment_free_config_id_seq RESTART;")

	sql := "INSERT INTO corvomq.environment_free_config (host, port, username, password, message_prefix) VALUES ($1, $2, $3, $4, $5);"
	_, _ = pgxpool.Prepare("insert", sql)
	defer pgxpool.Deallocate("insert")

	os.Remove("./gnatsd.conf")
	f, err := os.OpenFile("./gnatsd.conf", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	defer f.Close()

	s := "    {user: \"%v\", password: \"%v\", permissions: { publish = \"%v.*\", subscribe = \"%v.*\" }}\n"

	_, _ = f.WriteString("host: '0.0.0.0'\n")
	_, _ = f.WriteString("port: 4222\n")
	_, _ = f.WriteString("\n")
	_, _ = f.WriteString("authorization {\n")
	_, _ = f.WriteString("  users = [\n")

	startsWithAlpha, _ := regexp.Compile("^[a-zA-Z]")

	i := 1
	for i <= 100 {
		log.Println(i)

		username, _ := shortid.Generate()
		for !startsWithAlpha.MatchString(username) {
			username, _ = shortid.Generate()
		}

		password, _ := shortid.Generate()
		for !startsWithAlpha.MatchString(password) {
			password, _ = shortid.Generate()
		}

		messagePrefix, _ := shortid.Generate()

		_, _ = f.WriteString(fmt.Sprintf(s, username, password, messagePrefix, messagePrefix))

		_, err := pgxpool.Exec("insert", "free.corvomq.com", 4222, username, password, messagePrefix)
		if err != nil {
			log.Println(err.Error())
		}

		i++
	}
	_, _ = f.WriteString("  ]\n")
	_, _ = f.WriteString("}\n")
	_, _ = f.WriteString("\n")

	_, _ = f.WriteString("max_payload: 128\n")

}
