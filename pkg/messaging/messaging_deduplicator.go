package messaging

import (
	"context"
	"sync"
	"time"
)

type Deduplicator interface {
	IsDuplicate(ctx context.Context, message Message) (bool, error)
	MarkProcessed(ctx context.Context, message Message, ttlSeconds int) error
}

type InMemoryDeduplicator struct {
	mu    sync.Mutex
	store map[string]time.Time
}

func NewInMemoryDeduplicator() *InMemoryDeduplicator {
	return &InMemoryDeduplicator{
		store: make(map[string]time.Time),
	}
}

func (d *InMemoryDeduplicator) IsDuplicate(_ context.Context, message Message) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	expiry, exists := d.store[message.Identifier()]
	if !exists {
		return false, nil
	}

	if time.Now().After(expiry) {
		delete(d.store, message.Identifier())
		return false, nil
	}

	return true, nil
}

func (d *InMemoryDeduplicator) MarkProcessed(_ context.Context, message Message, ttlSeconds int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.store[message.Identifier()] = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	return nil
}
