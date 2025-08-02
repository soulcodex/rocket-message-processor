package rocketevents

import (
	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

var (
	ErrRocketEventBuilderNotFound = errutil.NewError("rocket event builder not found")
)

var rocketEventBuilders = map[string]RocketEvent{
	RocketExplodedType:       new(RocketExploded),
	RocketLaunchedType:       new(RocketLaunched),
	RocketMissionChangedType: new(RocketMissionChanged),
	RocketSpeedIncreasedType: new(RocketSpeedIncreased),
	RocketSpeedDecreasedType: new(RocketSpeedDecreased),
}

func ResolveRocketEvent(eventType string, raw *RocketEventRaw) (RocketEvent, error) {
	if rocketEvent, ok := rocketEventBuilders[eventType]; ok && raw != nil {
		if parseErr := rocketEvent.FromRawEvent(raw); parseErr != nil {
			return nil, errutil.NewError("failed to parse rocket event").Wrap(parseErr)
		}

		return rocketEvent, nil
	}

	return nil, ErrRocketEventBuilderNotFound
}
