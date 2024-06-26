package iam_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
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

func TestGetInstanceProfile(t *testing.T) {
	mockClient := new(mockIAMClient)
	service := iamHelpers.AWSService{IAMClient: mockClient}

	tests := []struct {
		name      string
		mockSetup func()
		want      *types.InstanceProfile
		wantErr   bool
	}{
		{
			name: "successful instance profile retrieval",
			mockSetup: func() {
				mockClient.On("GetInstanceProfile", mock.Anything, mock.Anything).Return(&iam.GetInstanceProfileOutput{
					InstanceProfile: &types.InstanceProfile{
						InstanceProfileName: aws.String("testProfile"),
					},
				}, nil).Once()
			},
			want: &types.InstanceProfile{
				InstanceProfileName: aws.String("testProfile"),
			},
			wantErr: false,
		},
		{
			name: "failure in instance profile retrieval",
			mockSetup: func() {
				mockClient.On("GetInstanceProfile", mock.Anything, mock.Anything).Return(nil, errors.New("failure in AWS service")).Once()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			got, err := service.GetInstanceProfile("testProfile")
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
