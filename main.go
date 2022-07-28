package main

import (
	"database/sql"
	"log"

	"github.com/kashshof/simplebank/api"
	db "github.com/kashshof/simplebank/db/sqlc"
	"github.com/kashshof/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	server := api.NewServer(db.NewStore(conn))
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
