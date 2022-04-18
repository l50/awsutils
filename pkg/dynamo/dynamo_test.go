package dynamo

import (
	"log"
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
