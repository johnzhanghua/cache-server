package cache

import (
	"context"
	"errors"
	"time"
)

// Cacher is middleware cache interface
// it defines Get/Set/Expire/Timeouts/TTL methods
type Cacher interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Expire(ctx context.Context, key string, expiration time.Duration) error

	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	TTL() time.Duration
}

var (
	// ErrorServiceNotReachable is the error indicates service not reachable
	ErrorServiceNotReachable = errors.New("cache service not reachable")
	// ErrorKeyNotExists is the error indicates key not exists in redis
	ErrorKeyNotExists = errors.New("key not exists in cache")
)
