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

// IdentityClientAPI represents the interface needed to make calls
// to the AWS IAM service.
//
// **Attributes:**
//
// GetInstanceProfile: Function to get the instance profile.
// CreateRole: Function to create a role.
// AttachRolePolicy: Function to attach a policy to a role.
// PutRolePolicy: Function to put a policy to a role.
// CreateInstanceProfile: Function to create an instance profile.
// AddRoleToInstanceProfile: Function to add a role to an instance profile.
// DetachRolePolicy: Function to detach a policy from a role.
// DeleteRolePolicy: Function to delete a policy from a role.
// RemoveRoleFromInstanceProfile: Function to remove a role from an instance profile.
// DeleteInstanceProfile: Function to delete an instance profile.
// DeleteRole: Function to delete a role.
type IdentityClientAPI interface {
	GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error)
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
	PutRolePolicy(ctx context.Context, params *iam.PutRolePolicyInput, optFns ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error)
	CreateInstanceProfile(ctx context.Context, params *iam.CreateInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.CreateInstanceProfileOutput, error)
	AddRoleToInstanceProfile(ctx context.Context, params *iam.AddRoleToInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.AddRoleToInstanceProfileOutput, error)
	DetachRolePolicy(ctx context.Context, params *iam.DetachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DetachRolePolicyOutput, error)
	DeleteRolePolicy(ctx context.Context, params *iam.DeleteRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DeleteRolePolicyOutput, error)
	RemoveRoleFromInstanceProfile(ctx context.Context, params *iam.RemoveRoleFromInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.RemoveRoleFromInstanceProfileOutput, error)
	DeleteInstanceProfile(ctx context.Context, params *iam.DeleteInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.DeleteInstanceProfileOutput, error)
	DeleteRole(ctx context.Context, params *iam.DeleteRoleInput, optFns ...func(*iam.Options)) (*iam.DeleteRoleOutput, error)
}

// AWSService represents the AWS services needed for the application.
//
// **Attributes:**
//
// STSClient: Client to make calls to the AWS STS service.
// IAMClient: Client to make calls to the AWS IAM service.
type AWSService struct {
	STSClient STSClientAPI
	IAMClient IdentityClientAPI
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

// CreateRole creates a new role with the given name and assume role policy.
//
// **Parameters:**
//
// roleName: The name of the role to create.
// assumeRolePolicy: The policy that the role will assume.
//
// **Returns:**
//
// *iam.CreateRoleOutput: A pointer to the CreateRoleOutput.
// error: An error if any issue occurs while trying to create the role.
func (s *AWSService) CreateRole(roleName, assumeRolePolicy string) (*iam.CreateRoleOutput, error) {
	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
	}

	result, err := s.IAMClient.CreateRole(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %v", err)
	}

	return result, nil
}

// AttachRolePolicy attaches a policy to a role.
//
// **Parameters:**
//
// roleName: The name of the role to attach the policy to.
// policyArn: The ARN of the policy to attach.
//
// **Returns:**
//
// *iam.AttachRolePolicyOutput: A pointer to the AttachRolePolicyOutput.
// error: An error if any issue occurs while trying to attach the policy to the role.
func (s *AWSService) AttachRolePolicy(roleName, policyArn string) (*iam.AttachRolePolicyOutput, error) {
	input := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyArn),
	}

	result, err := s.IAMClient.AttachRolePolicy(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to attach role policy: %v", err)
	}

	return result, nil
}

// PutRolePolicy updates the policy for a role.
//
// **Parameters:**
//
// roleName: The name of the role to update the policy for.
// policyName: The name of the policy to update.
// policyDocument: The policy document to update.
//
// **Returns:**
//
// *iam.PutRolePolicyOutput: A pointer to the PutRolePolicyOutput.
// error: An error if any issue occurs while trying to update the policy for the role.
func (s *AWSService) PutRolePolicy(roleName, policyName, policyDocument string) (*iam.PutRolePolicyOutput, error) {
	input := &iam.PutRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(policyDocument),
	}

	result, err := s.IAMClient.PutRolePolicy(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to put role policy: %v", err)
	}

	return result, nil
}

// CreateInstanceProfile creates a new instance profile with the given name.
//
// **Parameters:**
//
// profileName: The name of the instance profile to create.
//
// **Returns:**
//
// *iam.CreateInstanceProfileOutput: A pointer to the CreateInstanceProfileOutput.
// error: An error if any issue occurs while trying to create the instance profile.
func (s *AWSService) CreateInstanceProfile(profileName string) (*iam.CreateInstanceProfileOutput, error) {
	input := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}

	result, err := s.IAMClient.CreateInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance profile: %v", err)
	}

	return result, nil
}

// AddRoleToInstanceProfile adds a role to an instance profile.
//
// **Parameters:**
//
// profileName: The name of the instance profile to add the role to.
// roleName: The name of the role to add to the instance profile.
//
// **Returns:**
//
// *iam.AddRoleToInstanceProfileOutput: A pointer to the AddRoleToInstanceProfileOutput.
// error: An error if any issue occurs while trying to add the role to the instance profile.
func (s *AWSService) AddRoleToInstanceProfile(profileName, roleName string) (*iam.AddRoleToInstanceProfileOutput, error) {
	input := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}

	result, err := s.IAMClient.AddRoleToInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to add role to instance profile: %v", err)
	}

	return result, nil
}

// DetachRolePolicy detaches a policy from a role.
//
// **Parameters:**
//
// roleName: The name of the role to detach the policy from.
// policyArn: The ARN of the policy to detach.
//
// **Returns:**
//
// *iam.DetachRolePolicyOutput: A pointer to the DetachRolePolicyOutput.
// error: An error if any issue occurs while trying to detach the policy from the role.
func (s *AWSService) DetachRolePolicy(roleName, policyArn string) (*iam.DetachRolePolicyOutput, error) {
	input := &iam.DetachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyArn),
	}
	result, err := s.IAMClient.DetachRolePolicy(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to detach role policy: %v", err)
	}
	return result, nil
}

// DeleteRolePolicy deletes a policy from a role.
//
// **Parameters:**
//
// roleName: The name of the role to delete the policy from.
// policyName: The name of the policy to delete.
//
// **Returns:**
//
// *iam.DeleteRolePolicyOutput: A pointer to the DeleteRolePolicyOutput.
// error: An error if any issue occurs while trying to delete the policy from the role.
func (s *AWSService) DeleteRolePolicy(roleName, policyName string) (*iam.DeleteRolePolicyOutput, error) {
	input := &iam.DeleteRolePolicyInput{
		RoleName:   aws.String(roleName),
		PolicyName: aws.String(policyName),
	}
	result, err := s.IAMClient.DeleteRolePolicy(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete role policy: %v", err)
	}
	return result, nil
}

// RemoveRoleFromInstanceProfile removes a role from an instance profile.
//
// **Parameters:**
//
// profileName: The name of the instance profile to remove the role from.
// roleName: The name of the role to remove from the instance profile.
//
// **Returns:**
//
// *iam.RemoveRoleFromInstanceProfileOutput: A pointer to the RemoveRoleFromInstanceProfileOutput.
// error: An error if any issue occurs while trying to remove the role from the instance profile.
func (s *AWSService) RemoveRoleFromInstanceProfile(profileName, roleName string) (*iam.RemoveRoleFromInstanceProfileOutput, error) {
	input := &iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(roleName),
	}
	result, err := s.IAMClient.RemoveRoleFromInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to remove role from instance profile: %v", err)
	}
	return result, nil
}

// DeleteInstanceProfile deletes an instance profile.
//
// **Parameters:**
//
// profileName: The name of the instance profile to delete.
//
// **Returns:**
//
// *iam.DeleteInstanceProfileOutput: A pointer to the DeleteInstanceProfileOutput.
// error: An error if any issue occurs while trying to delete the instance profile.
func (s *AWSService) DeleteInstanceProfile(profileName string) (*iam.DeleteInstanceProfileOutput, error) {
	input := &iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}
	result, err := s.IAMClient.DeleteInstanceProfile(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete instance profile: %v", err)
	}
	return result, nil
}

// DeleteRole deletes a role.
//
// **Parameters:**
//
// roleName: The name of the role to delete.
//
// **Returns:**
//
// *iam.DeleteRoleOutput: A pointer to the DeleteRoleOutput.
// error: An error if any issue occurs while trying to delete the role.
func (s *AWSService) DeleteRole(roleName string) (*iam.DeleteRoleOutput, error) {
	input := &iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	}
	result, err := s.IAMClient.DeleteRole(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete role: %v", err)
	}
	return result, nil
}
