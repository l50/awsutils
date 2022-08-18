package secretsmanager

import (
	"log"
	"testing"
)

var (
	err      error
	smParams = Params{
		Name:        "TestSecret",
		Description: "Test Secret",
		Value:       "123456",
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
