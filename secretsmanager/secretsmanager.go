package secretsmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Connection contains all of the relevant
// information to maintain
// a Secrets Manager connection.
type Connection struct {
	Client  *secretsmanager.SecretsManager
	Session *session.Session
	Params  Params
}

// Params provides parameter
// options for Secrets Manager.
type Params struct {
	Name        string
	Description string
	Value       string
	Created     time.Time
	Modified    time.Time
}

// createClient is a helper function that
// returns a new secretsmanager session.
func createClient() (*secretsmanager.SecretsManager, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create secrets manager service client
	svc := secretsmanager.New(sess)

	return svc, sess
}

// CreateConnection creates a connection
// with Secrets Manager and returns it.
func CreateConnection() Connection {
	smConnection := Connection{}
	smConnection.Client, smConnection.Session = createClient()

	return smConnection
}

// CreateSecret creates an input `secretName`
// with the specified `secretValue`.
func CreateSecret(client *secretsmanager.SecretsManager,
	secretName string, secretDesc string, secretValue string) error {
	_, err := client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		Description:  aws.String(secretDesc),
		SecretString: aws.String(secretValue),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeLimitExceededException:
				return errors.New(secretsmanager.ErrCodeLimitExceededException)
			case secretsmanager.ErrCodeEncryptionFailure:
				return errors.New(secretsmanager.ErrCodeEncryptionFailure)
			case secretsmanager.ErrCodeResourceExistsException:
				return errors.New(secretsmanager.ErrCodeResourceExistsException)
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeMalformedPolicyDocumentException:
				return errors.New(secretsmanager.ErrCodeMalformedPolicyDocumentException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			case secretsmanager.ErrCodePreconditionNotMetException:
				return errors.New(secretsmanager.ErrCodePreconditionNotMetException)
			case secretsmanager.ErrCodeDecryptionFailure:
				return errors.New(secretsmanager.ErrCodeDecryptionFailure)
			default:
				return aerr
			}
		} else {
			return err
		}
	}
	return nil
}

// UpdateSecret updates an existing secret
func UpdateSecret(client *secretsmanager.SecretsManager, secretName string, secretValue string) error {
	_, err := client.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secretName),
		SecretString: aws.String(secretValue),
	})
	if err != nil {
		// Handle the error
		return fmt.Errorf("error updating secret: %v", err)
	}
	return nil
}

// CreateOrUpdateSecret creates a new secret or updates an existing one.
func CreateOrUpdateSecret(client *secretsmanager.SecretsManager, secretName string, secretDesc string, secretValue string) error {
	_, err := client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         aws.String(secretName),
		Description:  aws.String(secretDesc),
		SecretString: aws.String(secretValue),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceExistsException:
				// Secret already exists, update it.
				_, err := client.UpdateSecret(&secretsmanager.UpdateSecretInput{
					SecretId:     aws.String(secretName),
					SecretString: aws.String(secretValue),
				})
				if err != nil {
					return err
				}
			default:
				return aerr
			}
		} else {
			return err
		}
	}
	return nil
}

// DeleteSecret deletes an input `secretName`.
// It will attempt to do so forcefully if `forceDelete`
// is set to true.
func DeleteSecret(client *secretsmanager.SecretsManager, secretName string, forceDelete bool) error {
	_, err := client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretName),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeInvalidParameterException:
				return errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeInternalServiceError:
				return errors.New(secretsmanager.ErrCodeInternalServiceError)
			default:
				return aerr
			}
		} else {
			return err
		}
	}

	return nil
}

// GetSecret returns the value of an input `secretName`.
func GetSecret(client *secretsmanager.SecretsManager, secretName string) (string, error) {
	secret, err := client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return "", errors.New(secretsmanager.ErrCodeResourceNotFoundException)
			case secretsmanager.ErrCodeInvalidParameterException:
				return "", errors.New(secretsmanager.ErrCodeInvalidParameterException)
			case secretsmanager.ErrCodeInvalidRequestException:
				return "", errors.New(secretsmanager.ErrCodeInvalidRequestException)
			case secretsmanager.ErrCodeDecryptionFailure:
				return "", errors.New(secretsmanager.ErrCodeDecryptionFailure)
			case secretsmanager.ErrCodeInternalServiceError:
				return "", errors.New(secretsmanager.ErrCodeInternalServiceError)
			default:
				return "", aerr
			}
		} else {
			return "", err
		}
	}

	return *secret.SecretString, nil
}

// ReplicateSecret replicates a secret with the specified `secretName`
// to multiple target regions.
func ReplicateSecret(connection Connection, secretName string, newSecretName string, targetRegions []string) error {
	// Get the existing secret value
	secretValue, err := GetSecret(connection.Client, secretName)
	if err != nil {
		return fmt.Errorf("error getting secret value: %v", err)
	}

	// Replicate the secret to the target regions
	for _, targetRegion := range targetRegions {
		targetSession, err := session.NewSession(&aws.Config{
			Region: aws.String(targetRegion),
		})
		if err != nil {
			return fmt.Errorf("error creating session for region %s: %v", targetRegion, err)
		}

		targetClient := secretsmanager.New(targetSession)
		if err := CreateOrUpdateSecret(targetClient, newSecretName, "", secretValue); err != nil {
			return fmt.Errorf("error replicating secret to region %s: %v", targetRegion, err)
		}
	}

	return nil
}
