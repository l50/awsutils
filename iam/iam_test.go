package iam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	iamHelpers "github.com/l50/awsutils/iam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSTSClient struct {
	mock.Mock
}

func (m *mockSTSClient) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	args := m.Called(ctx, params)
	var output *sts.GetCallerIdentityOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*sts.GetCallerIdentityOutput)
	}
	return output, args.Error(1)
}

func TestGetAWSIdentity(t *testing.T) {
	mockClient := new(mockSTSClient)
	service := iamHelpers.AWSService{STSClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		want      *iamHelpers.AWSIdentity
		wantErr   bool
	}{
		{
			name: "successful AWS identity retrieval",
			mockSetup: func() {
				mockClient.On("GetCallerIdentity", mock.Anything, mock.Anything).Return(&sts.GetCallerIdentityOutput{
					Account: aws.String("123456789012"),
					Arn:     aws.String("arn:aws:iam::123456789012:user/TestUser"),
					UserId:  aws.String("TestUser"),
				}, nil).Once()
			},
			want: &iamHelpers.AWSIdentity{
				Account: "123456789012",
				ARN:     "arn:aws:iam::123456789012:user/TestUser",
				UserID:  "TestUser",
			},
			wantErr: false,
		},
		{
			name: "failure in AWS caller identity retrieval",
			mockSetup: func() {
				mockClient.On("GetCallerIdentity", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			got, err := service.GetAWSIdentity()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

type mockIAMClient struct {
	mock.Mock
}

func (m *mockIAMClient) GetInstanceProfile(ctx context.Context, params *iam.GetInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.GetInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.GetInstanceProfileOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.GetInstanceProfileOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.CreateRoleOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.CreateRoleOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.AttachRolePolicyOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.AttachRolePolicyOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) PutRolePolicy(ctx context.Context, params *iam.PutRolePolicyInput, optFns ...func(*iam.Options)) (*iam.PutRolePolicyOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.PutRolePolicyOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.PutRolePolicyOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) CreateInstanceProfile(ctx context.Context, params *iam.CreateInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.CreateInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.CreateInstanceProfileOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.CreateInstanceProfileOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) AddRoleToInstanceProfile(ctx context.Context, params *iam.AddRoleToInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.AddRoleToInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.AddRoleToInstanceProfileOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.AddRoleToInstanceProfileOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) DetachRolePolicy(ctx context.Context, params *iam.DetachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DetachRolePolicyOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.DetachRolePolicyOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.DetachRolePolicyOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) DeleteRolePolicy(ctx context.Context, params *iam.DeleteRolePolicyInput, optFns ...func(*iam.Options)) (*iam.DeleteRolePolicyOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.DeleteRolePolicyOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.DeleteRolePolicyOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) RemoveRoleFromInstanceProfile(ctx context.Context, params *iam.RemoveRoleFromInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.RemoveRoleFromInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.RemoveRoleFromInstanceProfileOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.RemoveRoleFromInstanceProfileOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) DeleteInstanceProfile(ctx context.Context, params *iam.DeleteInstanceProfileInput, optFns ...func(*iam.Options)) (*iam.DeleteInstanceProfileOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.DeleteInstanceProfileOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.DeleteInstanceProfileOutput)
	}
	return output, args.Error(1)
}

func (m *mockIAMClient) DeleteRole(ctx context.Context, params *iam.DeleteRoleInput, optFns ...func(*iam.Options)) (*iam.DeleteRoleOutput, error) {
	args := m.Called(ctx, params)
	var output *iam.DeleteRoleOutput
	if args.Get(0) != nil {
		output = args.Get(0).(*iam.DeleteRoleOutput)
	}
	return output, args.Error(1)
}

func TestCreateRole(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role creation",
			mockSetup: func() {
				mockClient.On("CreateRole", mock.Anything, mock.Anything).Return(&iam.CreateRoleOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role creation",
			mockSetup: func() {
				mockClient.On("CreateRole", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.CreateRole("testRole", "testPolicy")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestAttachRolePolicy(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful policy attachment",
			mockSetup: func() {
				mockClient.On("AttachRolePolicy", mock.Anything, mock.Anything).Return(&iam.AttachRolePolicyOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in policy attachment",
			mockSetup: func() {
				mockClient.On("AttachRolePolicy", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.AttachRolePolicy("testRole", "testPolicyArn")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestPutRolePolicy(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role policy put",
			mockSetup: func() {
				mockClient.On("PutRolePolicy", mock.Anything, mock.Anything).Return(&iam.PutRolePolicyOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role policy put",
			mockSetup: func() {
				mockClient.On("PutRolePolicy", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.PutRolePolicy("testRole", "testPolicy", "testPolicyDocument")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCreateInstanceProfile(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful instance profile creation",
			mockSetup: func() {
				mockClient.On("CreateInstanceProfile", mock.Anything, mock.Anything).Return(&iam.CreateInstanceProfileOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in instance profile creation",
			mockSetup: func() {
				mockClient.On("CreateInstanceProfile", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.CreateInstanceProfile("testInstanceProfile")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestAddRoleToInstanceProfile(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role addition to instance profile",
			mockSetup: func() {
				mockClient.On("AddRoleToInstanceProfile", mock.Anything, mock.Anything).Return(&iam.AddRoleToInstanceProfileOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role addition to instance profile",
			mockSetup: func() {
				mockClient.On("AddRoleToInstanceProfile", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.AddRoleToInstanceProfile("testInstanceProfile", "testRole")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestDetachRolePolicy(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful policy detachment",
			mockSetup: func() {
				mockClient.On("DetachRolePolicy", mock.Anything, mock.Anything).Return(&iam.DetachRolePolicyOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in policy detachment",
			mockSetup: func() {
				mockClient.On("DetachRolePolicy", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.DetachRolePolicy("testRole", "testPolicyArn")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestDeleteRolePolicy(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role policy deletion",
			mockSetup: func() {
				mockClient.On("DeleteRolePolicy", mock.Anything, mock.Anything).Return(&iam.DeleteRolePolicyOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role policy deletion",
			mockSetup: func() {
				mockClient.On("DeleteRolePolicy", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.DeleteRolePolicy("testRole", "testPolicy")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRemoveRoleFromInstanceProfile(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role removal from instance profile",
			mockSetup: func() {
				mockClient.On("RemoveRoleFromInstanceProfile", mock.Anything, mock.Anything).Return(&iam.RemoveRoleFromInstanceProfileOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role removal from instance profile",
			mockSetup: func() {
				mockClient.On("RemoveRoleFromInstanceProfile", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.RemoveRoleFromInstanceProfile("testInstanceProfile", "testRole")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestDeleteInstanceProfile(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful instance profile deletion",
			mockSetup: func() {
				mockClient.On("DeleteInstanceProfile", mock.Anything, mock.Anything).Return(&iam.DeleteInstanceProfileOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in instance profile deletion",
			mockSetup: func() {
				mockClient.On("DeleteInstanceProfile", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.DeleteInstanceProfile("testInstanceProfile")
			if tc.wantErr {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func TestDeleteRole(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful role deletion",
			mockSetup: func() {
				mockClient.On("DeleteRole", mock.Anything, mock.Anything).Return(&iam.DeleteRoleOutput{}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "failure in role deletion",
			mockSetup: func() {
				mockClient.On("DeleteRole", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			_, err := service.DeleteRole("testRole")
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}
