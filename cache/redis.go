package cache

import (
	"context"
	"time"
)

const (
	defaultTTL = 10 * time.Minute
)

// RedisCache implements Cacher interface
type RedisCache struct {
	client *RdbClient
}

// NewRedisCache create instance of RedisCache,
// and return as Cacher interface
func NewRedisCache() (Cacher, error) {
	rdbClient, err := NewRdbClient(nil)
	if err != nil {
		return nil, err
	}
	return &RedisCache{
		client: rdbClient,
	}, nil
}

// Get gets value from redis cache by key
func (r *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	return r.client.Get(ctx, key, value)
}

// Set sets key/value/expiration to redis
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration)
}

// Expire expires the key by expiration duration
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration)
}

// ReadTimeout returns ReadTimeout setting of cache
func (r *RedisCache) ReadTimeout() time.Duration {
	return r.client.ReadTimeout()
}

// WriteTimeout retures WriteTimeout setting of cache
func (r *RedisCache) WriteTimeout() time.Duration {
	return r.client.WriteTimeout()
}

// TTL returns the time-to-live for key in cache
func (r *RedisCache) TTL() time.Duration {
	return defaultTTL
}
