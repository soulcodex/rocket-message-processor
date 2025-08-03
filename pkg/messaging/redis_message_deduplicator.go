package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultTTLSeconds = 29
)

type RedisDeduplicator struct {
	client     *redis.Client
	messageTTL int
}

func NewDefaultRedisDeduplicator(client *redis.Client) *RedisDeduplicator {
	return &RedisDeduplicator{
		client:     client,
		messageTTL: defaultTTLSeconds,
	}
}

func (d *RedisDeduplicator) IsDuplicate(ctx context.Context, message Message) (bool, error) {
	exists, err := d.client.Exists(ctx, message.Identifier()).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check message duplicity: %w", err)
	}
	return exists == 1, nil
}

func (d *RedisDeduplicator) MarkProcessed(ctx context.Context, message Message) error {
	status := d.client.SetNX(ctx, message.Identifier(), "1", time.Duration(d.messageTTL)*time.Second)
	if _, err := status.Result(); err != nil {
		return fmt.Errorf("failed to mark message as processed: %w", err)
	}

	return nil
}
