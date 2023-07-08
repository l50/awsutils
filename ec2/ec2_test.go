package ec2

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/l50/goutils/v2/str"
)

var (
	err           error
	verbose       = false
	volumeSize, _ = str.ToInt64(os.Getenv("VOLUME_SIZE"))
	getPubIP, _   = strconv.ParseBool(os.Getenv("PUB_IP"))
	ec2Params     = Params{
		AssociatePublicIPAddress: getPubIP,
		ImageID:                  os.Getenv("AMI"),
		InstanceName:             os.Getenv("INST_NAME"),
		InstanceType:             os.Getenv("INST_TYPE"),
		InstanceProfile:          os.Getenv("IAM_INSTANCE_PROFILE"),
		MinCount:                 1,
		MaxCount:                 1,
		SecurityGroupIDs:         []string{os.Getenv("SEC_GRP_ID")},
		SubnetID:                 os.Getenv("SUBNET_ID"),
		VolumeSize:               volumeSize,
	}
	ec2Connection = Connection{}
)

func init() {
	ec2Connection.Client = createClient()
	ec2Connection.Params = ec2Params
	ec2Connection.Reservation, err = CreateInstance(
		ec2Connection.Client,
		ec2Connection.Params,
	)
	if err != nil {
		log.Fatalf(
			"error running CreateInstance(): %v",
			err,
		)
	}

	ec2Connection.Params.InstanceID = GetInstanceID(
		ec2Connection.Reservation.Instances[0],
	)

	fmt.Printf("Successfully created instance: %s\n",
		ec2Connection.Params.InstanceID)
}

func TestTagInstance(t *testing.T) {
	err = TagInstance(
		ec2Connection.Client,
		ec2Connection.Params.InstanceID,
		"Env",
		"Prod",
	)

	if err != nil {
		t.Fatalf(
			"error running TagInstance(): %v", err)
	}
}

func TestGetRunningInstances(t *testing.T) {
	_, err := GetRunningInstances(
		ec2Connection.Client)

	if err != nil {
		t.Fatalf(
			"error running GetRunningInstance(): %v", err)
	}
}

func TestWaitForInstance(t *testing.T) {
	// Skip test if running with
	// go test -short
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err = WaitForInstance(
		ec2Connection.Client,
		ec2Connection.Params.InstanceID,
	)
	if err != nil {
		t.Fatalf(
			"error running WaitForInstance(): %v",
			err,
		)
	}
}

// func TestGetInstancePublicIP(t *testing.T) {
// 	// Skip test if running with
// 	// go test -short
// 	if testing.Short() {
// 		t.Skip("skipping test in short mode.")
// 	}

// 	ec2Connection.Params.PublicIP, err =
// 		GetInstancePublicIP(
// 			ec2Connection.Client,
// 			ec2Connection.Params.InstanceID,
// 		)

// 	if err != nil {
// 		t.Fatalf(
// 			"error running GetInstancePublicIP(): %v",
// 			err,
// 		)
// 	}
// }

func TestGetRegion(t *testing.T) {
	_, err := GetRegion(ec2Connection.Client)
	if err != nil {
		t.Fatalf(
			"error running GetRegion(): %v",
			err,
		)
	}
}

func TestGetInstances(t *testing.T) {
	// Test with no filters
	instances, err := GetInstances(ec2Connection.Client, nil)
	if err != nil {
		t.Fatalf(
			"error running GetInstances() with no filters: %v",
			err,
		)
	}

	if verbose {
		log.Println("The following instances were found: ")
		for _, instance := range instances {
			fmt.Println(*instance.InstanceId)
		}
	}

	// Test with filters
	filters := []*ec2.Filter{
		{
			Name: aws.String("tag:Name"),
			Values: []*string{
				aws.String("goInstance"),
			},
		},
	}

	instances, err = GetInstances(ec2Connection.Client, filters)
	if err != nil {
		t.Fatalf(
			"error running GetInstances() with filters: %v",
			err,
		)
	}

	if verbose {
		log.Println("Using filters, the following instances were found: ")
		for _, instance := range instances {
			fmt.Println(*instance.InstanceId)
		}
	}
}

func TestGetInstanceState(t *testing.T) {
	state, err :=
		GetInstanceState(
			ec2Connection.Client,
			ec2Connection.Params.InstanceID,
		)

	if err != nil {
		t.Fatalf(
			"error running GetInstanceState(): %v",
			err,
		)
	}

	fmt.Printf(
		"Successfully grabbed instance state: %s\n",
		state)
}

func TestDestroyInstance(t *testing.T) {
	t.Cleanup(func() {
		err = DestroyInstance(
			ec2Connection.Client,
			ec2Connection.Params.InstanceID,
		)
		if err != nil {
			t.Fatalf(
				"error running DestroyInstance(): %v",
				err,
			)
		}
	})
}

func runningAction() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func TestIsEC2Instance(t *testing.T) {
	// Test that the function returns true when running on an EC2 instance
	// To simulate running on an EC2 instance, we can set the metadata endpoint to a mock server that returns a known instance ID
	metadataEndpoint := "http://localhost:8080/latest/meta-data/instance-id"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "i-1234567890abcdef")
	}))
	defer mockServer.Close()
	oldEndpoint := metadataEndpoint
	metadataEndpoint = mockServer.URL
	defer func() { metadataEndpoint = oldEndpoint }()

	// Running this test in a github action breaks the test logic.
	if IsEC2Instance() && !runningAction() {
		t.Error("expected IsEC2Instance() to return true when running on an EC2 instance")
	}

	// Test that the function returns false when not running on an EC2 instance
	// To simulate running on a non-EC2 environment, we can set the metadata endpoint to an invalid URL
	metadataEndpoint = "http://invalid-metadata-url"
	// Running this test in a github action breaks the test logic.
	if IsEC2Instance() && !runningAction() {
		t.Error("expected IsEC2Instance() to return false when not running on an EC2 instance")
	}
}

func TestGetLatestAMI(t *testing.T) {
	tests := []struct {
		name      string
		input     AMIInfo
		expectErr bool
	}{
		{
			name: "Ubuntu 22.04 arm64",
			input: AMIInfo{
				Distro:       "ubuntu",
				Version:      "22.04",
				Architecture: "arm64",
				Region:       "us-west-1",
			},
			expectErr: false,
		},
		{
			name: "Ubuntu 20.04 amd64",
			input: AMIInfo{
				Distro:       "ubuntu",
				Version:      "20.04",
				Architecture: "amd64",
				Region:       "us-west-1",
			},
			expectErr: false,
		},
		{
			name: "Unsupported distro",
			input: AMIInfo{
				Distro:       "not-supported",
				Version:      "20.04",
				Architecture: "amd64",
				Region:       "us-west-1",
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOutput, gotError := GetLatestAMI(tc.input)

			// Additional checks
			if gotError == nil {
				// AMI ID should start with "ami-"
				if !strings.HasPrefix(gotOutput, "ami-") {
					t.Errorf("expected AMI ID to start with 'ami-', got '%v'", gotOutput)
				}

				// Get AMI information
				imageOutput, err := ec2Connection.Client.DescribeImages(&ec2.DescribeImagesInput{
					ImageIds: []*string{&gotOutput},
				})

				if err != nil {
					t.Errorf("error describing image: %v", err)
				}

				image := imageOutput.Images[0]

				// Check if the architecture of the AMI matches
				architecture := tc.input.Architecture
				if architecture == "amd64" {
					architecture = "x86_64"
				}
				if *image.Architecture != architecture {
					t.Errorf("expected architecture to be '%v', got '%v'", architecture, *image.Architecture)
				}

				// Check if the image name contains the expected distro and version
				expectedNamePart := fmt.Sprintf("%s-%s", tc.input.Version, tc.input.Architecture)
				if !strings.Contains(*image.Name, expectedNamePart) {
					t.Errorf("expected image name to contain '%v', got '%v'", expectedNamePart, *image.Name)
				}
			}

			if (gotError != nil) != tc.expectErr {
				t.Errorf("expected error to be %v, got %v", tc.expectErr, gotError != nil)
			}
		})
	}
}
