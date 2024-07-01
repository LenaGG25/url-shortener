package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis client
type Redis struct {
	client     *redis.Client
	expiration time.Duration
}

// NewRedis constructor
func NewRedis(expiration time.Duration, opt *redis.Options) *Redis {
	return &Redis{
		redis.NewClient(opt),
		expiration,
	}
}

// Set value by key(shortURL)
func (r *Redis) Set(ctx context.Context, shortURL string, originalURL string) error {
	return r.client.Set(ctx, shortURL, originalURL, r.expiration).Err()
}

// Get value by key(shortURL)
func (r *Redis) Get(ctx context.Context, shortURL string) (string, error) {
	res := r.client.Get(ctx, shortURL)

	var originalURL string
	if err := res.Scan(&originalURL); err != nil {
		return "", ErrFindURL
	}

	return originalURL, nil
}
