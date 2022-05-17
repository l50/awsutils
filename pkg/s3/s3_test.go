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
	s3Connection.Client = createClient()
	s3Connection.Params = s3Params
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	err = CreateBucket(s3Connection.Client,
		s3Connection.Params.BucketName)
	if err != nil {
		log.Fatalf(
			"error running CreateBucket(): %v",
			err,
		)
	}

	log.Println(
		"Waiting for test bucket to finish " +
			"initialization - please wait",
	)

	err = WaitForBucket(
		s3Connection.Client,
		s3Connection.Params.BucketName,
	)
	if err != nil {
		log.Fatalf(
			"error running WaitForBucket(): %v",
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

func TestDestroyBucket(t *testing.T) {
	t.Cleanup(func() {
		err = DestroyBucket(s3Connection.Client,
			s3Connection.Params.BucketName)
		if err != nil {
			t.Fatalf(
				"error running DestroyBucket(): %v",
				err,
			)
		}
	})
}
