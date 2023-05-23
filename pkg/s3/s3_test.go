package s3

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	goutils "github.com/l50/goutils"
)

var (
	err        error
	randStr, _ = goutils.RandomString(10)
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
	verbose = false
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
	uploadFile := "testFile"
	if goutils.CreateEmptyFile(uploadFile) {
		if err := goutils.AppendToFile(uploadFile, "teststring123"); err == nil {
			if err := UploadBucketFile(s3Connection.Session, randStr, uploadFile); err != nil {
				t.Fatalf(
					"error uploading %s to %s: %v",
					uploadFile,
					randStr,
					err,
				)
			}
		}
	}

	if err := os.Remove(uploadFile); err != nil {
		t.Fatalf("failed to remove %s: %v", uploadFile, err)
	}

	downloadedFileName, err := DownloadBucketFile(s3Connection.Session, randStr, uploadFile, "/tmp/"+uploadFile)
	if err != nil {
		t.Fatalf(
			"error downloading %s from %s: %v",
			uploadFile,
			randStr,
			err,
		)
	}

	downloadedFile, err := os.Open(downloadedFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer downloadedFile.Close()

	downloadedFileContent, err := io.ReadAll(downloadedFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(downloadedFileContent) == 0 {
		t.Fatalf(
			"downloaded file %s is empty", downloadedFileName,
		)
	}

	if err := os.Remove(downloadedFileName); err != nil {
		t.Fatalf("failed to remove %s: %v", downloadedFileName, err)
	}
}

func TestDestroyBucket(t *testing.T) {
	t.Cleanup(func() {
		if err := EmptyBucket(s3Connection.Client, s3Params.BucketName); err != nil {
			t.Fatalf(
				"error running EmptyBucket(): %v",
				err,
			)
		}

		if err := DestroyBucket(s3Connection.Client, s3Params.BucketName); err != nil {
			t.Fatalf(
				"error running DestroyBucket(): %v",
				err,
			)
		}
	})
}
