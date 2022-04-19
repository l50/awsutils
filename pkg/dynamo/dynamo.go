package dynamo

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// Connection contains all of the
// relevant information to maintain
// a connection with DynamoDB.
type Connection struct {
	Client *dynamodb.DynamoDB
	Params Params
}

// Params provides parameter
// options for a DynamoDB table.
type Params struct {
	ID        uuid.UUID
	TableName string
	Created   time.Time
	Modified  time.Time
}

// createClient is a helper function that
// returns a new dynamo session.
func createClient() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return svc
}

// CreateConnection creates a connection
// with DynamoDB and returns a Connection.
func CreateConnection() Connection {
	dynamoConnection := Connection{}
	dynamoConnection.Client = createClient()

	return dynamoConnection
}

// GetRegion returns the region associated with the input
// dynamo client.
func GetRegion(client *dynamodb.DynamoDB) (string, error) {
	region := client.Config.Region

	if region == nil {
		return "", errors.New("failed to retrieve region")
	}

	return *region, nil
}

// GetTables returns all dynamoDB tables the input client
// has access to.
// Resource:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-list-tables.html
func GetTables(client *dynamodb.DynamoDB) ([]*string, error) {
	var result = &dynamodb.ListTablesOutput{}
	var err error
	input := &dynamodb.ListTablesInput{}

	for {
		result, err = client.ListTables(input)
		if err != nil {
			return result.TableNames, err
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return result.TableNames, nil
}

// CreateTable creates a table with the input
// dynamoConnection.
func CreateTable(dynamoConnection Connection) error {
	_, err :=
		dynamoConnection.Client.CreateTable(
			&dynamodb.CreateTableInput{
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("id"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("HASH"),
					},
				},
				TableName:   aws.String(dynamoConnection.Params.TableName),
				BillingMode: aws.String("PAY_PER_REQUEST"),
			})

	if err != nil {
		return err
	}

	return nil
}
