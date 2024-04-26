package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STSClientAPI represents the interface needed to make calls
// to the AWS STS service.
//
// **Attributes:**
//
// GetCallerIdentity: Function to get the caller identity.
type STSClientAPI interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// IAMClientAPI represents the interface needed to make calls
// to the AWS IAM service.
//
// **Attributes:**
//
// GetInstanceProfile: Function to get the instance profile.
type IAMClientAPI interface {
	GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)
}

// AWSService represents the AWS services needed for the application.
//
// **Attributes:**
//
// STSClient: Client to make calls to the AWS STS service.
// IAMClient: Client to make calls to the AWS IAM service.
type AWSService struct {
	STSClient STSClientAPI
	IAMClient IAMClientAPI
}

// AWSIdentity represents the identity of an AWS account.
//
// **Attributes:**
//
// Account: AWS account ID.
// ARN: AWS ARN associated with the account.
// UserID: AWS user ID associated with the account.
type AWSIdentity struct {
	Account string
	ARN     string
	UserID  string
}

// NewAWSService creates a new AWSService with the default AWS configuration.
//
// **Returns:**
//
// *AWSService: A pointer to the newly created AWSService.
// error: An error if any issue occurs while trying to create the AWSService.
func NewAWSService() (*AWSService, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	iamClient := iam.NewFromConfig(cfg)
	return &AWSService{STSClient: stsClient, IAMClient: iamClient}, nil
}

// GetAWSIdentity retrieves the AWS identity of the caller.
//
// **Returns:**
//
// *AWSIdentity: A pointer to the AWSIdentity of the caller.
// error: An error if any issue occurs while trying to get the AWS identity.
func (s *AWSService) GetAWSIdentity() (*AWSIdentity, error) {
	result, err := s.STSClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS caller identity: %v", err)
	}
	identity := &AWSIdentity{
		Account: *result.Account,
		ARN:     *result.Arn,
		UserID:  *result.UserId,
	}
	return identity, nil
}

// GetInstanceProfile retrieves the instance profile for a given profile name.
//
// **Parameters:**
//
// profileName: The name of the profile to retrieve.
//
// **Returns:**
//
// *types.InstanceProfile: A pointer to the InstanceProfile.
// error: An error if any issue occurs while trying to get the instance profile.
func (s *AWSService) GetInstanceProfile(profileName string) (*types.InstanceProfile, error) {
	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}

	result, err := s.IAMClient.GetInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result.InstanceProfile, nil
}
