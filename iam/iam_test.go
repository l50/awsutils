package iam_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/github.com/l50/aws/iam"
)

func TestGetIamInstanceProfile(t *testing.T) {
	tests := []struct {
		name           string
		profileName    string
		expectError    bool
	}{
		{
			name:        "valid instance profile",
			profileName: "validProfile", // replace with an existing IAM instance profile name
			expectError: false,
		},
		{
			name:        "invalid instance profile",
			profileName: "invalidProfile", // replace with a non-existing IAM instance profile name
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := iam.GetIamInstanceProfile(tc.profileName)
			if (err != nil) != tc.expectError {
				t.Errorf("GetIamInstanceProfile(%v) returned error %v, expectError %v", tc.profileName, err, tc.expectError)
			}
		})
	}
}

