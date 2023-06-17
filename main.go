package main

import (
	"database/sql"
	"log"

	"github.com/binbomb/goapp/simplebank/api"
	db "github.com/binbomb/goapp/simplebank/db/sqlc"
	"github.com/binbomb/goapp/simplebank/utils"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".") // app.env
	if err != nil {
		log.Fatal("cannot load config to file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server ", err)
	}
	server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server ", err)
	}
}
