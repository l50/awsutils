package s3_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/l50/awsutils/s3"
)

func ExampleGetBuckets() {
	s3Connection := s3.CreateConnection()

	buckets, _ := s3.GetBuckets(s3Connection.Client)
	for _, b := range buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}
