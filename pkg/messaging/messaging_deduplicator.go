package messaging

import (
	"context"
)

type Deduplicator interface {
	IsDuplicate(ctx context.Context, message Message) (bool, error)
	MarkProcessed(ctx context.Context, message Message) error
}
