package s3

import (
	"log"
	"testing"
	"time"

	utils "github.com/l50/goutils"
)

var (
	err        error
	randStr, _ = utils.RandomString(10)
	s3Params   = Params{
		BucketName: randStr,
		Created:    time.Now(),
		Modified:   time.Now(),
	}
	s3Connection = Connection{}
	verbose      bool
)

func init() {
	verbose = false
	s3Connection.Client, s3Connection.Session = createClient()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	err = CreateBucket(s3Connection.Client,
		randStr)
	if err != nil {
		log.Fatalf(
			"error running CreateBucket(): %v",
			err,
		)
	}
}

func TestGetBuckets(t *testing.T) {
	result, err := GetBuckets(s3Connection.Client)
	if err != nil {
		t.Fatalf(
			"error running GetBuckets(): %v",
			err,
		)
	}
	if verbose {
		for _, n := range result {
			log.Println("Bucket: ", *n.Name)
		}
	}
}

func TestDownloadBucketFile(t *testing.T) {
	file := "testFile"
	if utils.CreateEmptyFile(file) {
		err := UploadBucketFile(s3Connection.Session, randStr, file)
		if err != nil {
			t.Fatalf(
				"error uploading %s to %s: %v",
				file,
				randStr,
				err,
			)
		}
	}

	err = DownloadBucketFile(s3Connection.Session, randStr, file)
	if err != nil {
		t.Fatalf(
			"error downloading %s from %s: %v",
			file,
			randStr,
			err,
		)
	}
}

func TestDestroyBucket(t *testing.T) {
	t.Cleanup(func() {
		err = EmptyBucket(s3Connection.Client, s3Params.BucketName)
		if err != nil {
			t.Fatalf(
				"error running EmptyBucket(): %v",
				err,
			)
		}

		err = DestroyBucket(s3Connection.Client,
			s3Params.BucketName)
		if err != nil {
			t.Fatalf(
				"error running DestroyBucket(): %v",
				err,
			)
		}
	})
}
