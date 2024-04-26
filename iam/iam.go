package iam

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws"
)

// STSClientAPI defines the methods of the STS client that are used.
type STSClientAPI interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

// IdentityService provides AWS identity services.
type IdentityService struct {
	Client STSClientAPI
}

// AWSIdentity holds the identity information for the AWS caller.
type AWSIdentity struct {
	Account string
	ARN     string
	UserID  string
}

// IAMClientAPI defines the methods of the IAM client that are used.
type IAMClientAPI interface {
	GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)
}

// InstanceProfileService provides AWS instance profile services.
type InstanceProfileService struct {
	Client IAMClientAPI
}

// NewIdentityService creates a new IdentityService with an STS client.
func NewIdentityService() (*IdentityService, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %v", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	return &IdentityService{Client: stsClient}, nil
}

func (s *IdentityService) GetAWSIdentity() (*AWSIdentity, error) {
	result, err := s.Client.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
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

// NewInstanceProfileService creates a new InstanceProfileService with an IAM client.
func NewInstanceProfileService() (*InstanceProfileService, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	iamClient := iam.NewFromConfig(cfg)
	return &InstanceProfileService{Client: iamClient}, nil
}

func (s *InstanceProfileService) GetInstanceProfile(profileName string) (*types.InstanceProfile, error) {
	input := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}

	result, err := s.Client.GetInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result.InstanceProfile, nil
}
