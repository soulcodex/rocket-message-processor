package rocketqueries

import (
	"context"
	"fmt"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type SearchRocketsQuery struct {
	Sort string
	Asc  bool
}

func (q *SearchRocketsQuery) Type() string {
	return "search_rockets_query"
}

type SearchRocketsQueryHandler struct {
	repository rocketdomain.RocketRepository
}

func NewSearchRocketsQueryHandler(repository rocketdomain.RocketRepository) *SearchRocketsQueryHandler {
	return &SearchRocketsQueryHandler{
		repository: repository,
	}
}

func (h *SearchRocketsQueryHandler) Handle(ctx context.Context, q *SearchRocketsQuery) (RocketsResponse, error) {
	rockets, err := h.repository.Search(ctx, q.Sort, q.Asc)
	if err != nil {
		return nil, fmt.Errorf("failed to search rockets: %w", err)
	}

	return newRocketsResponseFromPrimitives(rockets.Primitives()), nil
}
