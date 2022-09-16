package main

import (
	"database/sql"
	"log"
	"simplebank/cmd/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:8080/simple_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("could not start server")
	}
}
