package main

import (
	"database/sql"
	"log"
	"smutaxi/api"
	db "smutaxi/db/sqlc"
	"smutaxi/util"

	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgres://root:secret@localhost:5432/smutaxi?sslmode=disable"
// 	serverAddress = "0.0.0.0:3000"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot read config file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot connect to server: ", err)
	}
}
