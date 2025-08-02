package rocketevents

import (
	"context"
	"fmt"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type DeleteRocketOnRocketExploded struct {
	updater *rocketdomain.RocketUpdater
}

func NewDeleteRocketOnRocketExploded(updater *rocketdomain.RocketUpdater) *DeleteRocketOnRocketExploded {
	return &DeleteRocketOnRocketExploded{
		updater: updater,
	}
}

func (e *DeleteRocketOnRocketExploded) Handle(ctx context.Context, evt *RocketExploded) (interface{}, error) {
	_, err := e.updater.Update(ctx, evt.RocketID, rocketdomain.WithSoftDeletion(evt.OccurredOn))
	if err != nil {
		return nil, fmt.Errorf("error updating rocket: %w", err)
	}

	return struct{}{}, nil
}
