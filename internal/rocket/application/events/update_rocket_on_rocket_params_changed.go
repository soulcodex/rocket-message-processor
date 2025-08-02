package rocketevents

import (
	"context"
	"fmt"

	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
)

type UpdateRocketOnRocketParamsChanged struct {
	updater *rocketdomain.RocketUpdater
}

func NewUpdateRocketOnRocketParamsChanged(updater *rocketdomain.RocketUpdater) *UpdateRocketOnRocketParamsChanged {
	return &UpdateRocketOnRocketParamsChanged{
		updater: updater,
	}
}

func (e *UpdateRocketOnRocketParamsChanged) Handle(ctx context.Context, evt RocketEvent) (interface{}, error) {
	switch event := evt.(type) {
	case *RocketMissionChanged:
		return e.handleMissionChanged(ctx, event)
	case *RocketSpeedIncreased:
		return e.handleSpeedIncreased(ctx, event)
	case *RocketSpeedDecreased:
		return e.handleSpeedDecreased(ctx, event)
	default:
		return nil, fmt.Errorf("unhandled rocket params change event type: %s", evt.Type())
	}
}

func (e *UpdateRocketOnRocketParamsChanged) handleSpeedIncreased(ctx context.Context, evt *RocketSpeedIncreased) (interface{}, error) {
	_, err := e.updater.Update(ctx, evt.RocketID, rocketdomain.WithLaunchSpeed(int64(evt.Amount), evt.OccurredOn))
	if err != nil {
		return nil, fmt.Errorf("error updating rocket: %w", err)
	}

	return struct{}{}, nil
}

func (e *UpdateRocketOnRocketParamsChanged) handleSpeedDecreased(ctx context.Context, evt *RocketSpeedDecreased) (interface{}, error) {
	_, err := e.updater.Update(ctx, evt.RocketID, rocketdomain.WithLaunchSpeed(int64(-evt.Amount), evt.OccurredOn))
	if err != nil {
		return nil, fmt.Errorf("error updating rocket: %w", err)
	}

	return struct{}{}, nil
}

func (e *UpdateRocketOnRocketParamsChanged) handleMissionChanged(ctx context.Context, evt *RocketMissionChanged) (interface{}, error) {
	_, err := e.updater.Update(ctx, evt.RocketID, rocketdomain.WithMission(evt.NewMission, evt.OccurredOn))
	if err != nil {
		return nil, fmt.Errorf("error updating rocket: %w", err)
	}

	return struct{}{}, nil
}
