package eventbus

import (
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
)

type EventBus = bus.Bus

func InitEventBus() EventBus {
	return bus.InitSyncBus()
}
