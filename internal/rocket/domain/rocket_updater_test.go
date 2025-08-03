package rocketdomain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	rockettest "github.com/soulcodex/rockets-message-processor/test/rocket"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rocketmock "github.com/soulcodex/rockets-message-processor/internal/rocket/domain/mock"
)

func TestRocketUpdater_Update(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name          string
		id            string
		updates       []rocketdomain.RocketUpdaterFunc
		setupMocks    func(repo *rocketmock.RocketRepositoryMock)
		validation    func(t *testing.T, rocket *rocketdomain.Rocket)
		expectedError string
	}{
		{
			name:    "should update rocket mission successfully",
			id:      "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
			updates: []rocketdomain.RocketUpdaterFunc{rocketdomain.WithMission("Moon landing", now)},
			setupMocks: func(repo *rocketmock.RocketRepositoryMock) {
				repo.FindFunc = func(_ context.Context, _ rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
					return rockettest.NewRocketMother().Build(t), nil
				}
				repo.SaveFunc = func(_ context.Context, r *rocketdomain.Rocket) error {
					return nil
				}
			},
		},
		{
			name: "should fail on invalid rocket ID",
			id:   "invalid-id",
			setupMocks: func(_ *rocketmock.RocketRepositoryMock) {
				// not needed
			},
			expectedError: "invalid rocket id",
		},
		{
			name:    "should fail when rocket not found",
			id:      "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
			updates: []rocketdomain.RocketUpdaterFunc{rocketdomain.WithMission("X", now)},
			setupMocks: func(repo *rocketmock.RocketRepositoryMock) {
				repo.FindFunc = func(_ context.Context, _ rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
					return nil, errors.New("not found")
				}
			},
			expectedError: "failed to find rocket",
		},
		{
			name: "should decrease rocket speed with delta",
			id:   "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
			updates: []rocketdomain.RocketUpdaterFunc{
				rocketdomain.WithLaunchSpeedDelta(int64(-1000), time.Now().Add(30*time.Minute)),
			},
			validation: func(t *testing.T, rocket *rocketdomain.Rocket) {
				assert.Equal(t, int64(4000), rocket.Primitives().LaunchSpeed)
			},
			setupMocks: func(repo *rocketmock.RocketRepositoryMock) {
				repo.FindFunc = func(_ context.Context, _ rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
					return rockettest.NewRocketMother(
						rockettest.WithLaunchSpeed(5000),
					).Build(t), nil
				}
				repo.SaveFunc = func(_ context.Context, r *rocketdomain.Rocket) error {
					assert.Equal(t, int64(4000), r.Primitives().LaunchSpeed)
					return nil
				}
			},
		},
		{
			name:    "should fail when repository save fails",
			id:      "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
			updates: []rocketdomain.RocketUpdaterFunc{rocketdomain.WithMission("New mission", now)},
			setupMocks: func(repo *rocketmock.RocketRepositoryMock) {
				repo.FindFunc = func(_ context.Context, _ rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
					return rockettest.NewRocketMother().Build(t), nil
				}
				repo.SaveFunc = func(_ context.Context, _ *rocketdomain.Rocket) error {
					return errors.New("db error")
				}
			},
			expectedError: "failed to update rocket: db error",
		},
		{
			name:    "should soft delete rocket",
			id:      "122c31b0-a3c4-411a-bc07-5f342f0d78e4",
			updates: []rocketdomain.RocketUpdaterFunc{rocketdomain.WithSoftDeletion(time.Now().Add(30 * time.Minute))},
			validation: func(t *testing.T, rocket *rocketdomain.Rocket) {
				assert.NotNil(t, rocket.Primitives().DeletedAt)
			},
			setupMocks: func(repo *rocketmock.RocketRepositoryMock) {
				repo.FindFunc = func(_ context.Context, _ rocketdomain.RocketID) (*rocketdomain.Rocket, error) {
					return rockettest.NewRocketMother().Build(t), nil
				}
				repo.SaveFunc = func(_ context.Context, _ *rocketdomain.Rocket) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &rocketmock.RocketRepositoryMock{}
			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			updater := rocketdomain.NewRocketUpdater(repo)
			rocket, err := updater.Update(ctx, tt.id, tt.updates...)

			if tt.expectedError == "" {
				require.NoError(t, err)
				require.NotNil(t, rocket)

				if tt.validation != nil {
					tt.validation(t, rocket)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, rocket)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
