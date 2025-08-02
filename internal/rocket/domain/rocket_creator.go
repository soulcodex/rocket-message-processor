package rocketdomain

import (
	"context"
	"fmt"
	"time"
)

type RocketCreateParams struct {
	ID          string
	RocketType  string
	LaunchSpeed int64
	Mission     string
	At          time.Time
}
type RocketCreator struct {
	repository RocketRepository
}

func NewRocketCreator(repository RocketRepository) *RocketCreator {
	return &RocketCreator{
		repository: repository,
	}
}

func (r *RocketCreator) Create(ctx context.Context, dto RocketCreateParams) (*Rocket, error) {
	id, err := NewRocketID(dto.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid rocket id: %w", err)
	}

	rocketType, err := NewRocketType(dto.RocketType)
	if err != nil {
		return nil, fmt.Errorf("invalid rocket type: %w", err)
	}

	launchSpeed, err := NewLaunchSpeed(dto.LaunchSpeed)
	if err != nil {
		return nil, fmt.Errorf("invalid rocket launch speed: %w", err)
	}

	mission, err := NewMission(dto.Mission)
	if err != nil {
		return nil, fmt.Errorf("invalid mission: %w", err)
	}

	rocket := NewRocket(id, rocketType, launchSpeed, mission, dto.At)

	if saveErr := r.repository.Save(ctx, rocket); saveErr != nil {
		return nil, fmt.Errorf("failed to save rocket: %w", saveErr)
	}

	return rocket, nil
}
