package messaging

import (
	"sync"
	"time"

	"github.com/soulcodex/rockets-message-processor/pkg/bus"
)

var _ bus.Dto = (*BaseMessage)(nil)

type BaseMessage struct {
	msgID       string
	msgType     string
	msgMetadata []byte
	msgData     []byte
	msgTime     time.Time
	ackOnce     sync.Once
	nackOnce    sync.Once
	ack         chan struct{}
	noAck       chan struct{}
}

func (bm *BaseMessage) Type() string {
	return bm.msgType
}

func (bm *BaseMessage) Time() time.Time {
	return bm.msgTime
}

func (bm *BaseMessage) DataRaw() []byte {
	return bm.msgData
}

func (bm *BaseMessage) Metadata() []byte {
	return bm.msgMetadata
}

func NewBaseMessage(
	msgID string,
	msgType string,
	msgMetadata []byte,
	msgData []byte,
	msgTime time.Time,
) (*BaseMessage, error) {
	return &BaseMessage{
		msgID:       msgID,
		msgType:     msgType,
		msgMetadata: msgMetadata,
		msgData:     msgData,
		msgTime:     msgTime,
		ackOnce:     sync.Once{},
		nackOnce:    sync.Once{},
		ack:         make(chan struct{}),
		noAck:       make(chan struct{}),
	}, nil
}

// Ack marks the message as acknowledged
func (bm *BaseMessage) Ack() {
	bm.ackOnce.Do(func() {
		select {
		case <-bm.ack:
			// Channel is already closed
		default:
			close(bm.ack)
		}
	})
}

// Nack marks the message as not acknowledged
func (bm *BaseMessage) Nack() {
	bm.nackOnce.Do(func() {
		select {
		case <-bm.noAck:
			// Channel is already closed
		default:
			close(bm.noAck)
		}
	})
}

// Acked returns a channel that is closed when the message is acknowledged
func (bm *BaseMessage) Acked() <-chan struct{} {
	return bm.ack
}

// Nacked returns a channel that is closed when the message is not acknowledged
func (bm *BaseMessage) Nacked() <-chan struct{} {
	return bm.noAck
}
