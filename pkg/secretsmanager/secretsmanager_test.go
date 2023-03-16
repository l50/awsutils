package secretsmanager

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

func TestReplicateSecret(t *testing.T) {
	newSecretName, _ := utils.RandomString(10)
	targetRegions := []string{"us-west-1", "eu-west-1"}

	if err := ReplicateSecret(smConnection, smParams.Name, newSecretName, targetRegions); err != nil {
		t.Fatalf("error replicating secret: %v", err)
	}

	// Verify that the new secret exists in each target region
	for _, region := range targetRegions {
		targetSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			t.Fatalf("error creating session for region %s: %v", region, err)
		}

		targetClient := secretsmanager.New(targetSession)

		if secret, err := GetSecret(targetClient, newSecretName); err != nil {
			t.Fatalf("error getting replicated secret in region %s: %v", region, err)
		}

		// Cleanup: delete the replicated secret in each target region
		if err := DeleteSecret(targetClient, newSecretName, true); err != nil {
			t.Fatalf("error deleting replicated secret in region %s: %v", region, err)
		}
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
