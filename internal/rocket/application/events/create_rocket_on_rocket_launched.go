package rocketevents

import (
	"context"
	"fmt"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type CreateRocketOnRocketLaunched struct {
	creator *rocketdomain.RocketCreator
}

func NewCreateRocketOnRocketLaunched(creator *rocketdomain.RocketCreator) *CreateRocketOnRocketLaunched {
	return &CreateRocketOnRocketLaunched{
		creator: creator,
	}
}

func (e *CreateRocketOnRocketLaunched) Handle(ctx context.Context, evt *RocketLaunched) (interface{}, error) {
	input := rocketdomain.RocketCreateParams{
		ID:          evt.RocketID,
		RocketType:  evt.RocketType,
		LaunchSpeed: evt.LaunchSpeed,
		Mission:     evt.Mission,
		At:          evt.OccurredOn,
	}

	_, err := e.creator.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error creating rocket: %w", err)
	}

	return struct{}{}, nil
}
