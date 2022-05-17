package s3

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Connection contains all of the relevant
// information to maintain
// an S3 connection.
type Connection struct {
	Client *s3.S3
	Params Params
}

// Params provides parameter
// options for an S3 bucket.
type Params struct {
	BucketName string
	Created    time.Time
	Modified   time.Time
}

// createClient is a helper function that
// returns a new s3 session.
func createClient() *s3.S3 {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create S3 service client
	svc := s3.New(sess)

	return svc
}

// CreateConnection creates a connection
// with DynamoDB and returns it.
func CreateConnection() Connection {
	s3Connection := Connection{}
	s3Connection.Client = createClient()

	return s3Connection
}

// WaitForBucket waits for the input bucket to
// finish being created.
func WaitForBucket(client *s3.S3, bucket string) error {
	err := client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return err
	}

	return nil
}

// CreateBucket creates a bucket with the input
// `bucketName`.
func CreateBucket(client *s3.S3, bucketName string) error {
	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	return nil
}

// GetBuckets returns all s3 buckets
// that the input client has access to.
func GetBuckets(client *s3.S3) ([]*s3.Bucket, error) {
	var err error
	input := &s3.ListBucketsInput{}

	result, err := client.ListBuckets(input)
	if err != nil {
		return result.Buckets, err
	}

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	return result.Buckets, nil
}

// DestroyBucket destroys a bucket with the input
// `bucketName`.
func DestroyBucket(client *s3.S3, bucketName string) error {

	_, err := client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	err = client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return err
	}

	return nil
}
