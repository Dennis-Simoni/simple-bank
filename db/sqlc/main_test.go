package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	c, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("could not read config in test")
	}
	testDB, err = sql.Open(c.DBDriver, c.DBSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
