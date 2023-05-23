package s3

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Connection contains all of the relevant
// information to maintain
// an S3 connection.
type Connection struct {
	Client  *s3.S3
	Session *session.Session
	Params  Params
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
func createClient() (*s3.S3, *session.Session) {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	// Create S3 service client
	svc := s3.New(sess)

	return svc, sess
}

// CreateConnection creates a connection
// with DynamoDB and returns it.
func CreateConnection() Connection {
	s3Connection := Connection{}
	s3Connection.Client, s3Connection.Session = createClient()

	return s3Connection
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

// EmptyBucket deletes everything found in the input `bucketName`.
func EmptyBucket(client *s3.S3, bucketName string) error {

	iter := s3manager.NewDeleteListIterator(client, &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	// Iterate through bucket and delete each discovered object.
	if err := s3manager.NewBatchDeleteWithClient(client).Delete(aws.BackgroundContext(), iter); err != nil {
		return err
	}

	return nil
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

// UploadBucketFile uploads a file found at the file path specified with (`uploadFP`)
// to the input `bucketName`.
func UploadBucketFile(sess *session.Session, bucketName string, uploadFP string) error {
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(uploadFP)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uploadFP),
		Body:   file,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %q to %q\n", uploadFP, bucketName)

	return nil
}

// DownloadBucketFile downloads a file found at the file path specified with
// the input objectKey to the input bucketName.
func DownloadBucketFile(sess *session.Session, bucketName string, objectKey string, downloadFP string) (string, error) {
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(downloadFP)
	if err != nil {
		return "", err
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully downloaded", file.Name(), numBytes, "bytes")

	return file.Name(), nil
}
