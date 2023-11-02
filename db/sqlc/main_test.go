// To connect to the db and communicate with it before writing an unit test to test about other function such as CreateAccount, DeleteAccount.
package db

import (
	"database/sql"
	"log"
	"os"
	"simple_bank/util"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
