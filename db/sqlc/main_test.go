package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/scipiia/snippetbox/util"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:secret@localhost:5432/snippetbox?sslmode=disable"
// )

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := util.LiadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
