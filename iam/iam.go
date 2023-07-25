package iam

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// GetIamInstanceProfile retrieves the IAM instance profile by its name
func GetIamInstanceProfile(profileName string) (*iam.InstanceProfile, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := iam.New(sess)
	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: &profileName,
	}

	result, err := svc.GetInstanceProfile(input)
	if err != nil {
		return nil, err
	}

	return result.InstanceProfile, nil
}
