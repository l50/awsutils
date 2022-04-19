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
		ID:        uuid.New(),
		TableName: "testTable",
		Created:   time.Now(),
		Modified:  time.Now(),
	}
	dbConnection = Connection{}
)

func init() {
	dbConnection.Client = createClient()
	dbConnection.Params = dbParams
	if err != nil {
		log.Fatalf(
			"error running CreateClient(): %v",
			err,
		)
	}

	err = CreateTable(dbConnection.Client,
		dbConnection.Params.TableName)
	if err != nil {
		log.Fatalf(
			"error running CreateTable(): %v",
			err,
		)
	}

	log.Println(
		"Waiting for test table to finish initialization - please wait",
	)

	err = WaitForTable(
		dbConnection.Client,
		dbConnection.Params.TableName,
	)
	if err != nil {
		log.Fatalf(
			"error running WaitForTable(): %v",
			err,
		)
	}
}

func TestGetRegion(t *testing.T) {
	_, err = GetRegion(dbConnection.Client)
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
			"error running GetTables(): %v",
			err,
		)
	}

	for _, n := range result {
		log.Println("Table: ", *n)
	}
}

func TestDestroyTable(t *testing.T) {
	t.Cleanup(func() {
		err = DestroyTable(dbConnection.Client,
			dbConnection.Params.TableName)
		if err != nil {
			t.Fatalf(
				"error running DestroyTable(): %v",
				err,
			)
		}
	})
}
