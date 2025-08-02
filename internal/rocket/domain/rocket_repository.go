package rocketdomain

import (
	"context"
)

//go:generate moq -pkg rocketdomainmock -out mock/rocket_repository_moq.go . RocketRepository
type RocketRepository interface {
	Find(ctx context.Context, id RocketID) (*Rocket, error)
	Save(ctx context.Context, r *Rocket) error
}
