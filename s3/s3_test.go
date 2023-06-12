package s3_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/l50/awsutils/s3"
	fileutils "github.com/l50/goutils/v2/file"
	"github.com/l50/goutils/v2/str"
)

var (
	err          error
	randStr, _   = str.GenRandom(10)
	s3Connection = s3.Connection{}
	verbose      bool
)

func init() {
	verbose = false
	s3Connection = s3.CreateConnection()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	if err := s3.CreateBucket(s3Connection.Client,
		randStr); err != nil {
		log.Fatalf(
			"error running CreateBucket(): %v",
			err,
		)
	}
}

func TestS3Functions(t *testing.T) {
	tests := []struct {
		name        string
		bucketName  string
		uploadFile  string
		downloadDir string
	}{
		{
			name:        "test 1",
			bucketName:  randStr,
			uploadFile:  "testFile",
			downloadDir: "/tmp/",
		},
	}

	s3Connection := s3.CreateConnection()
	if err := s3.CreateBucket(s3Connection.Client, randStr); err != nil {
		t.Fatalf("error running CreateBucket(): %v", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// TestGetBuckets
			if _, err := s3.GetBuckets(s3Connection.Client); err != nil {
				t.Fatalf("error running GetBuckets(): %v", err)
			}

			// TestDownloadBucketFile
			testContent := []byte("teststring123")
			if err := fileutils.Create(tc.uploadFile, testContent); err != nil {
				t.Fatalf("error creating %s: %v", tc.uploadFile, err)
			}

			if err := fileutils.Append(tc.uploadFile, string(testContent)); err != nil {
				t.Fatalf("error appending to %s: %v", tc.uploadFile, err)
			} else {
				if err := s3.UploadBucketFile(s3Connection.Session, tc.bucketName, tc.uploadFile); err != nil {
					t.Fatalf("error uploading %s to %s: %v", tc.uploadFile, tc.bucketName, err)
				}
			}

			if err := os.Remove(tc.uploadFile); err != nil {
				t.Fatalf("failed to remove %s: %v", tc.uploadFile, err)
			}

			downloadPath := filepath.Join(tc.downloadDir, tc.uploadFile)
			_, err = s3.DownloadBucketFile(s3Connection.Session, tc.bucketName, tc.uploadFile, downloadPath)
			if err != nil {
				t.Fatalf("error downloading %s from %s: %v", tc.uploadFile, tc.bucketName, err)
			}

			file, err := os.Open(downloadPath)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			content, err := io.ReadAll(file)
			if err != nil {
				t.Fatal(err)
			}

			if len(content) == 0 {
				t.Fatalf("downloaded file %s is empty", downloadPath)
			}

			if err := os.Remove(downloadPath); err != nil {
				t.Fatalf("failed to remove %s: %v", downloadPath, err)
			}

			// TestDestroyBucket
			if err := s3.EmptyBucket(s3Connection.Client, tc.bucketName); err != nil {
				t.Fatalf("error running EmptyBucket(): %v", err)
			}

			if err := s3.DestroyBucket(s3Connection.Client, tc.bucketName); err != nil {
				t.Fatalf("error running DestroyBucket(): %v", err)
			}
		})
	}
}
