package ec2

import (
	"log"
)

func ExampleGetLatestAMI() {
	info := AMIInfo{
		Distro:       "ubuntu",
		Version:      "20.04",
		Architecture: "amd64",
		Region:       "us-west-2",
	}

	amiID, err := GetLatestAMI(info)
	if err != nil {
		log.Fatalf("failed to get latest AMI: %v", err)
	}

	log.Println("Latest AMI ID:", amiID)
}
