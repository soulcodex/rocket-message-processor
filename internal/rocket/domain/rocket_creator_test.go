package rocketdomain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rocketmock "github.com/soulcodex/rockets-message-processor/internal/rocket/domain/mock"
)

func validRocketCreateParams() rocketdomain.RocketCreateParams {
	return rocketdomain.RocketCreateParams{
		ID:          "8c696c9c-eb1c-422a-85d4-fd4f7c9fa676",
		RocketType:  "Falcon9",
		LaunchSpeed: 27000,
		Mission:     "Satellite Deployment",
		At:          time.Now(),
	}
}

func TestRocketCreator_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		input         rocketdomain.RocketCreateParams
		setupMock     func(repo *rocketmock.RocketRepositoryMock)
		expectedError string
	}{
		{
			name:  "should create rocket successfully",
			input: validRocketCreateParams(),
			setupMock: func(repo *rocketmock.RocketRepositoryMock) {
				repo.SaveFunc = func(ctx context.Context, r *rocketdomain.Rocket) error {
					return nil
				}
			},
			expectedError: "",
		},
		{
			name: "should fail on invalid rocket ID",
			input: func() rocketdomain.RocketCreateParams {
				in := validRocketCreateParams()
				in.ID = ""
				return in
			}(),
			setupMock:     func(repo *rocketmock.RocketRepositoryMock) {},
			expectedError: "invalid rocket id",
		},
		{
			name: "should fail on invalid rocket type",
			input: func() rocketdomain.RocketCreateParams {
				in := validRocketCreateParams()
				in.RocketType = ""
				return in
			}(),
			setupMock:     func(repo *rocketmock.RocketRepositoryMock) {},
			expectedError: "invalid rocket type",
		},
		{
			name: "should fail on invalid mission",
			input: func() rocketdomain.RocketCreateParams {
				in := validRocketCreateParams()
				in.Mission = ""
				return in
			}(),
			setupMock:     func(repo *rocketmock.RocketRepositoryMock) {},
			expectedError: "invalid mission",
		},
		{
			name:  "should fail when repository save fails",
			input: validRocketCreateParams(),
			setupMock: func(repo *rocketmock.RocketRepositoryMock) {
				repo.SaveFunc = func(ctx context.Context, r *rocketdomain.Rocket) error {
					return errors.New("db error")
				}
			},
			expectedError: "failed to save rocket: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &rocketmock.RocketRepositoryMock{}
			tt.setupMock(mockRepo)

			creator := rocketdomain.NewRocketCreator(mockRepo)
			rocket, err := creator.Create(ctx, tt.input)

			if tt.expectedError == "" {
				require.NoError(t, err)
				require.NotNil(t, rocket)
			} else {
				require.Error(t, err)
				require.Nil(t, rocket)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
