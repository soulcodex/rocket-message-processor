package rocketdomain

import (
	"context"
	"fmt"
	"time"
)

type RocketUpdaterFunc func(rocket *Rocket) error

func WithLaunchSpeed(launchSpeed uint64, at time.Time) RocketUpdaterFunc {
	return func(rocket *Rocket) error {
		speed, err := NewLaunchSpeed(launchSpeed)
		if err != nil {
			return fmt.Errorf("invalid rocket launch speed: %w", err)
		}
		rocket.IncreaseLaunchSpeed(speed, at)
		return nil
	}
}

func WithMission(mission string, at time.Time) RocketUpdaterFunc {
	return func(rocket *Rocket) error {
		newMission, err := NewMission(mission)
		if err != nil {
			return fmt.Errorf("invalid mission: %w", err)
		}
		rocket.ChangeMission(newMission, at)
		return nil
	}
}

type RocketUpdater struct {
	repository RocketRepository
}

func NewRocketUpdater(repository RocketRepository) *RocketUpdater {
	return &RocketUpdater{
		repository: repository,
	}
}

func (r *RocketUpdater) Update(ctx context.Context, rocketID string, updates ...RocketUpdaterFunc) (*Rocket, error) {
	id, err := NewRocketID(rocketID)
	if err != nil {
		return nil, fmt.Errorf("invalid rocket id: %w", err)
	}

	rocket, err := r.repository.Find(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find rocket: %w", err)
	}

	for _, update := range updates {
		if updateErr := update(rocket); updateErr != nil {
			return nil, fmt.Errorf("failed to apply rocket update: %w", updateErr)
		}
	}

	if saveErr := r.repository.Save(ctx, rocket); saveErr != nil {
		return nil, fmt.Errorf("failed to update rocket: %w", saveErr)
	}

	return rocket, nil
}
