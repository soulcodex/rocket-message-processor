package messaging

import (
	"time"
)

// Message is an interface that defines the methods to access message properties.
type Message interface {
	Time() time.Time
	Schema() string
	DataRaw() []byte
	Metadata() []byte
}

// AckAwaiter is an interface that exposes the methods to wait for an ack or nack.
type AckAwaiter interface {
	Acked() <-chan struct{}
	Nacked() <-chan struct{}
}

// Acknowledgeable exposes the methods to acknowledge or reject a message.
type Acknowledgeable interface {
	Ack()
	Nack()
}
