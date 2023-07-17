package main

import (
	"database/sql"
	"log"

	"github.com/scipiia/snippetbox/api"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/util"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver     = "postgres"
// 	dbSource     = "postgresql://root:secret@localhost:5432/snippetbox?sslmode=disable"
// 	serverAdress = "0.0.0.0:8080"
// )

func main() {
	config, err := util.LiadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	query := db.New(conn)
	server := api.NewServer(query)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
