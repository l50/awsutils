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

	_, err := PutParam(ssmConnection.Client,
		&ssmParams.Name, &ssmParams.Value,
		&ssmParams.Type, &ssmParams.Overwrite)
	if err != nil {
		log.Fatalf(
			"error running CreateSSMParam(): %v",
			err,
		)
	}
}
func TestGetParam(t *testing.T) {
	result, err := GetParam(ssmConnection.Client,
		&ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
	fmt.Println(*result.Parameter.Value)
}

func TestDeleteParam(t *testing.T) {
	_, err := DeleteParam(ssmConnection.Client, &ssmParams.Name)
	if err != nil {
		t.Fatalf(
			"error running GetParam(): %v",
			err,
		)
	}
}
