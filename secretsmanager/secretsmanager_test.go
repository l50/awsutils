package secretsmanager

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/l50/goutils/v2/str"
)

var (
	err        error
	randStr, _ = str.GenRandom(10)
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
	newSecretName, _ := str.GenRandom(10)
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

		if _, err := GetSecret(targetClient, newSecretName); err != nil {
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

func TestUpdateSecret(t *testing.T) {
	// First, let's create a secret for testing update functionality
	testName, _ := str.GenRandom(10)
	testValue := "InitialValue"
	err := CreateSecret(smConnection.Client, testName, "Test Secret", testValue)
	if err != nil {
		t.Fatalf("error creating secret for TestUpdateSecret: %v", err)
	}

	// Now, let's update this secret
	newValue := "UpdatedValue"
	err = UpdateSecret(smConnection.Client, testName, newValue)
	if err != nil {
		t.Fatalf("error updating secret: %v", err)
	}

	// Retrieve the updated secret and check if it's updated correctly
	secretValue, err := GetSecret(smConnection.Client, testName)
	if err != nil {
		t.Fatalf("error getting updated secret: %v", err)
	}

	if secretValue != newValue {
		t.Fatalf("error in TestUpdateSecret, expected: %s, got: %s", newValue, secretValue)
	}

	// Cleanup: delete the secret we used for testing
	err = DeleteSecret(smConnection.Client, testName, true)
	if err != nil {
		t.Fatalf("error deleting secret for TestUpdateSecret: %v", err)
	}
}

func TestCreateOrUpdateSecret(t *testing.T) {
	testName, _ := str.GenRandom(10)
	testValue := "InitialValue"

	// Use CreateOrUpdateSecret to create a new secret
	err := CreateOrUpdateSecret(smConnection.Client, testName, "Test Secret", testValue)
	if err != nil {
		t.Fatalf("error creating secret for TestCreateOrUpdateSecret: %v", err)
	}

	// Verify if the secret was created successfully
	secretValue, err := GetSecret(smConnection.Client, testName)
	if err != nil {
		t.Fatalf("error getting created secret: %v", err)
	}

	if secretValue != testValue {
		t.Fatalf("error in TestCreateOrUpdateSecret, expected: %s, got: %s", testValue, secretValue)
	}

	// Now, let's update this secret using the same function
	newValue := "UpdatedValue"
	err = CreateOrUpdateSecret(smConnection.Client, testName, "Test Secret", newValue)
	if err != nil {
		t.Fatalf("error updating secret for TestCreateOrUpdateSecret: %v", err)
	}

	// Verify if the secret was updated successfully
	secretValue, err = GetSecret(smConnection.Client, testName)
	if err != nil {
		t.Fatalf("error getting updated secret: %v", err)
	}

	if secretValue != newValue {
		t.Fatalf("error in TestCreateOrUpdateSecret, expected: %s, got: %s", newValue, secretValue)
	}

	// Cleanup: delete the secret we used for testing
	err = DeleteSecret(smConnection.Client, testName, true)
	if err != nil {
		t.Fatalf("error deleting secret for TestCreateOrUpdateSecret: %v", err)
	}
}
