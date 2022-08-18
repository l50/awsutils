package ssm

import (
	"fmt"
	"log"
	"testing"
)

var (
	err       error
	ssmParams = Param{
		Name:      "TestParam",
		Value:     "123456",
		Type:      "String",
		Overwrite: true,
	}
	ssmConnection = Connection{}
	verbose       bool
)

func init() {
	verbose = false
	ssmConnection.Client, ssmConnection.Session = createClient()
	if err != nil {
		log.Fatalf(
			"error running createClient(): %v",
			err,
		)
	}

	err := PutParam(ssmConnection.Client,
		ssmParams.Name, ssmParams.Value,
		ssmParams.Type, ssmParams.Overwrite)
	if err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
			err,
		)
	}
}
func TestGetParam(t *testing.T) {
	result, err := GetParam(ssmConnection.Client,
		ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
	fmt.Println(result)
}

func TestDeleteParam(t *testing.T) {
	err := DeleteParam(ssmConnection.Client, ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
}

// TODO: create an instance as part of this test
// so that we can run this test.
// func TestRunCommand(t *testing.T) {
// 	instanceID := "TODO"
// 	command := []string{
// 		"whoami",
// 	}

// 	result, err := RunCommand(ssmConnection.Client, instanceID, command)
// 	if err != nil {
// 		t.Fatalf(
// 			"error running RunCommand(): %v",
// 			err,
// 		)
// 	}
// fmt.Println(result)
// }
