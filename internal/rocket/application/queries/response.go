package rocketqueries

import (
	"iter"
	"time"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type RocketsResponse []RocketResponse

func newRocketsResponseFromPrimitives(rockets iter.Seq[rocketdomain.RocketPrimitives]) RocketsResponse {
	response := make(RocketsResponse, 0)
	for rocket := range rockets {
		response = append(response, newRocketResponseFromPrimitives(rocket))
	}
	return response
}

type RocketResponse struct {
	ID          string
	RocketType  string
	LaunchSpeed int64
	Mission     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newRocketResponseFromPrimitives(p rocketdomain.RocketPrimitives) RocketResponse {
	return RocketResponse{
		ID:          p.ID,
		RocketType:  p.RocketType,
		LaunchSpeed: p.LaunchSpeed,
		Mission:     p.Mission,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
