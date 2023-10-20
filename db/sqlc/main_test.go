// To connect to the db and communicate with it before writing an unit test to test about other function such as CreateAccount, DeleteAccount.
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@127.0.0.1:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
