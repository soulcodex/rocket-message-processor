package rocketdomain

import (
	"context"
)

//go:generate moq -pkg rocketdomainmock -out mock/rocket_repository_moq.go . RocketRepository
type RocketRepository interface {
	Find(ctx context.Context, id RocketID) (*Rocket, error)
	Search(ctx context.Context, sortBy string, asc bool) (RocketCollection, error)
	Save(ctx context.Context, r *Rocket) error
}
