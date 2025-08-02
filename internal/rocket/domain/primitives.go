package rocketdomain

import (
	"time"
)

type RocketPrimitives struct {
	ID          string
	RocketType  string
	LaunchSpeed int64
	Mission     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func primitivesFromDomain(r *Rocket) RocketPrimitives {
	return RocketPrimitives{
		ID:          r.id.String(),
		RocketType:  r.rocketType.String(),
		LaunchSpeed: r.launchSpeed.Value(),
		Mission:     r.mission.String(),
		CreatedAt:   r.createdAt,
		UpdatedAt:   r.updatedAt,
		DeletedAt:   r.deletedAt,
	}
}
