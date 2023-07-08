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

func ExampleIsEC2Instance() {
	isEC2 := IsEC2Instance()
	if isEC2 {
		log.Println("Running on an EC2 instance")
	} else {
		log.Println("Not running on an EC2 instance")
	}
}
