package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
var testQueries *Queries
var testDb *sql.DB
*/
var testStore Store

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..") // app.env
	if err != nil {
		log.Fatal("cannot load config to file: ", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
