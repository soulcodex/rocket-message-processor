package rocketqueries

import (
	"context"
	"fmt"
	"time"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type FindRocketByIDResponse struct {
	ID          string
	RocketType  string
	LaunchSpeed int64
	Mission     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func newFindRocketByIDResponseFromPrimitives(p rocketdomain.RocketPrimitives) FindRocketByIDResponse {
	return FindRocketByIDResponse{
		ID:          p.ID,
		RocketType:  p.RocketType,
		LaunchSpeed: p.LaunchSpeed,
		Mission:     p.Mission,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

type FindRocketByIDQuery struct {
	RocketID string
}

func (q *FindRocketByIDQuery) Type() string {
	return "find_rocket_by_id_query"
}

type FindRocketByIDQueryHandler struct {
	repository rocketdomain.RocketRepository
}

func NewFindRocketByIDQueryHandler(repository rocketdomain.RocketRepository) *FindRocketByIDQueryHandler {
	return &FindRocketByIDQueryHandler{
		repository: repository,
	}
}

func (h *FindRocketByIDQueryHandler) Handle(ctx context.Context, q *FindRocketByIDQuery) (FindRocketByIDResponse, error) {
	rocketID, err := rocketdomain.NewRocketID(q.RocketID)
	if err != nil {
		return FindRocketByIDResponse{}, fmt.Errorf("invalid rocket ID provided: %w", err)
	}

	rocket, err := h.repository.Find(ctx, rocketID)
	if err != nil {
		return FindRocketByIDResponse{}, fmt.Errorf("error while finding rocket by ID: %w", err)
	}

	return newFindRocketByIDResponseFromPrimitives(rocket.Primitives()), nil
}
