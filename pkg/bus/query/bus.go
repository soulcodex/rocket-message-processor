package querybus

import (
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
)

type Bus = bus.Bus

func InitQueryBus() Bus {
	return bus.InitSyncBus()
}
