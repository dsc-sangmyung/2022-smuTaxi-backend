package main

import (
	"database/sql"
	"log"
	"smutaxi/api"
	db "smutaxi/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/smutaxi?sslmode=disable"
	serverAddress = "0.0.0.0:3000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot connect to server: ", err)
	}
}
