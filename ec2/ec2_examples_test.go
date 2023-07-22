package ec2_test

import (
	"fmt"
	"log"

	ec2utils "github.com/l50/awsutils/ec2"
)

func ExampleConnection_GetLatestAMI() {
	c := ec2utils.NewConnection()
	info := ec2utils.AMIInfo{
		Distro:       "ubuntu",
		Version:      "20.04",
		Architecture: "amd64",
		Region:       "us-west-1",
	}

	amiID, err := c.GetLatestAMI(info)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(amiID)
}

func ExampleIsEC2Instance() {
	isEC2 := ec2utils.IsEC2Instance()
	if isEC2 {
		log.Println("Running on an EC2 instance")
	} else {
		log.Println("Not running on an EC2 instance")
	}
}
