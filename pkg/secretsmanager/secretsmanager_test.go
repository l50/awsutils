package secretsmanager

import (
	"log"
	"testing"
	"time"

	utils "github.com/l50/goutils"
)

var (
	err        error
	randStr, _ = utils.RandomString(10)
	smParams   = Params{
		Name:        randStr,
		Description: "Test Secret",
		Value:       "123456",
		Created:     time.Now(),
		Modified:    time.Now(),
	}
	smConnection = Connection{}
	verbose      bool
)

func init() {
	verbose = false
	smConnection.Client, smConnection.Session = createClient()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	err := CreateSecret(smConnection.Client,
		smParams.Name, smParams.Description,
		smParams.Value)
	if err != nil {
		log.Fatalf(
			"error running CreateSecret(): %v",
			err,
		)
	}
}

func TestGetSecret(t *testing.T) {
	_, err := GetSecret(smConnection.Client,
		smParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetSecret(): %v",
			err,
		)
	}
}

func TestDeleteSecret(t *testing.T) {
	err := DeleteSecret(smConnection.Client,
		smParams.Name, true)
	if err != nil {
		t.Fatalf(
			"error running DeleteSecret(): %v",
			err,
		)
	}
}
