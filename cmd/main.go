package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)

func main() {

	c, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load configurations")
	}

	conn, err := sql.Open(c.DBDriver, c.DBSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(c.ServerAddress)
	if err != nil {
		log.Fatal("could not start server")
	}
}
