package rocketdomain

import (
	"context"
)

type RocketRepository interface {
	Exists(ctx context.Context, id RocketID) (bool, error)
	Find(ctx context.Context, id RocketID) (*Rocket, error)
	Save(ctx context.Context, r *Rocket) error
}
