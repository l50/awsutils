package s3

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Connection is a struct that contains all of the relevant
// information to maintain an S3 connection.
//
// Client: An AWS S3 client that can be used to interact with an S3 service.
// Session: An AWS session that is used to create the S3 client.
// Params: Parameters that are used when interacting with the S3 service.
type Connection struct {
	Client  *s3.S3
	Session *session.Session
	Params  Params
}

// Params is a struct that provides parameter
// options for an S3 bucket.
//
// BucketName: The name of the bucket.
// Created: The time the bucket was created.
// Modified: The last modified time of the bucket.
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
// with S3 and returns it.
//
// Returns:
// An s3.Connection struct containing the AWS S3 client and AWS session.
func CreateConnection() Connection {
	s3Connection := Connection{}
	s3Connection.Client, s3Connection.Session = createClient()

	return s3Connection
}

// CreateBucket creates a bucket with the input
// bucketName.
//
// Parameters:
// client: An AWS S3 client.
// bucketName: The name of the bucket to create.
//
// Returns:
// error: An error if the bucket could not be created.
func CreateBucket(client *s3.S3, bucketName string) error {
	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Printf("Bucket %s already exists and is owned by you\n", bucketName)
				return nil
			case s3.ErrCodeBucketAlreadyExists:
				return fmt.Errorf("bucket %s already exists", bucketName)
			default:
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// GetBuckets returns all s3 buckets
// that the input client has access to.
//
// Parameters:
// client: An AWS S3 client.
//
// Returns:
// []*s3.Bucket: A slice of pointers to AWS S3 bucket structs that the client has access to.
// error: An error if the buckets could not be listed.
func GetBuckets(client *s3.S3) ([]*s3.Bucket, error) {
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

// EmptyBucket deletes everything found in the input bucketName.
//
// Parameters:
// client: An AWS S3 client.
// bucketName: The name of the bucket to empty.
//
// Returns:
// error: An error if the bucket could not be emptied.
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
// bucketName.
//
// Parameters:
// client: An AWS S3 client.
// bucketName: The name of the bucket to destroy.
//
// Returns:
// error: An error if the bucket could not be destroyed.
func DestroyBucket(client *s3.S3, bucketName string) error {
	if _, err := client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	}); err != nil {
		return err
	}

	if err := client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}); err != nil {
		return err
	}

	return nil
}

// UploadBucketDir uploads a directory specified by dirPath
// to the bucket specified by bucketName.
//
// Parameters:
// sess: An AWS session.
// bucketName: The name of the bucket to upload to.
// dirPath: The file path of the directory to upload.
//
// Returns:
// error: An error if the directory could not be uploaded.
func UploadBucketDir(sess *session.Session, bucketName string, dirPath string) error {
	if err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err := UploadBucketFile(sess, bucketName, path)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded directory %q to %q\n", dirPath, bucketName)
	return nil
}

// UploadBucketFile uploads a file found at the file path specified with (uploadFP)
// to the input bucketName.
//
// Parameters:
// sess: An AWS session.
// bucketName: The name of the bucket to upload to.
// uploadFP: The file path of the file to upload.
//
// Returns:
// error: An error if the file could not be uploaded.
func UploadBucketFile(sess *session.Session, bucketName string, uploadFP string) error {
	uploader := s3manager.NewUploader(sess)

	if _, err := os.Stat(uploadFP); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(uploadFP)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uploadFP),
		Body:   file,
	}); err != nil {
		return err
	}

	fmt.Printf("Successfully uploaded %v to %v\n", uploadFP, bucketName)

	return nil
}

// DownloadBucketFile downloads a file found at the object key specified with
// the input objectKey from the bucket specified by bucketName, and writes it to downloadFP.
//
// Parameters:
// sess: An AWS session.
// bucketName: The name of the bucket to download from.
// objectKey: The key of the object to download.
// downloadFP: The file path to write the downloaded file to.
//
// Returns:
// string: The name of the downloaded file.
// error: An error if the file could not be downloaded.
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
