package rocketqueries

import (
	"context"
	"fmt"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

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

func (h *FindRocketByIDQueryHandler) Handle(ctx context.Context, q *FindRocketByIDQuery) (RocketResponse, error) {
	rocketID, err := rocketdomain.NewRocketID(q.RocketID)
	if err != nil {
		return RocketResponse{}, fmt.Errorf("invalid rocket ID provided: %w", err)
	}

	rocket, err := h.repository.Find(ctx, rocketID)
	if err != nil {
		return RocketResponse{}, fmt.Errorf("error while finding rocket by ID: %w", err)
	}

	return newRocketResponseFromPrimitives(rocket.Primitives()), nil
}
