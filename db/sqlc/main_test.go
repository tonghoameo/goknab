package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/binbomb/goapp/simplebank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..") // app.env
	if err != nil {
		log.Fatal("cannot load config to file: ", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db: ", err)
	}
	testQueries = New(testDb)
	os.Exit(m.Run())
}
