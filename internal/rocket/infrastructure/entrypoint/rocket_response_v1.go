package rocketentrypoint

import (
	"time"

	rocketqueries "github.com/soulcodex/rockets-message-processor/internal/rocket/application/queries"
)

type RocketsResponseV1 []*RocketResponseV1

func newRocketsResponseV1(rockets rocketqueries.RocketsResponse) RocketsResponseV1 {
	items := make(RocketsResponseV1, len(rockets))

	for i, r := range rockets {
		rocket := newRocketResponseV1(r)
		items[i] = &rocket
	}

	return items
}

type RocketResponseV1 struct {
	ID          string    `json:"id"`
	RocketType  string    `json:"rocket_type"`
	LaunchSpeed int64     `json:"launch_speed"`
	Mission     string    `json:"mission"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newRocketResponseV1(r rocketqueries.RocketResponse) RocketResponseV1 {
	return RocketResponseV1{
		ID:          r.ID,
		RocketType:  r.RocketType,
		LaunchSpeed: r.LaunchSpeed,
		Mission:     r.Mission,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
