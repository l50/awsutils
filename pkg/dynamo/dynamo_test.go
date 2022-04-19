package dynamo

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	err      error
	dbParams = Params{
		ID:       uuid.New(),
		Created:  time.Now(),
		Modified: time.Now(),
	}
	dbConnection = Connection{}
)

func init() {
	dbConnection.Client = createClient()
	dbConnection.Params = dbParams
	if err != nil {
		log.Fatalf(
			"error running CreateInstance(): %v",
			err,
		)
	}
}

func TestGetRegion(t *testing.T) {
	_, err := GetRegion(dbConnection.Client)
	if err != nil {
		t.Fatalf(
			"error running GetRegion(): %v",
			err,
		)
	}
}

func TestGetTables(t *testing.T) {
	result, err := GetTables(dbConnection.Client)
	if err != nil {
		t.Fatalf(
			"error running ListTables(): %v",
			err,
		)
	}

	for _, n := range result {
		log.Println("Table: ", *n)
	}
}
