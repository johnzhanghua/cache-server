package cache

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RdbClient is the wrapper of redis client
type RdbClient struct {
	client *redis.Client
}

const (
	defaultAddr     = "127.0.0.1:6379"
	defaultUser     = ""
	defaultPassword = ""
)

func defaultOptions() *redis.Options {
	options := &redis.Options{
		Addr:         defaultAddr,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Username:     defaultUser,
		Password:     defaultPassword,
		//TLSConfig:    &tls.Config{}, // TODO, enable tls configs
	}

	addr := os.Getenv("REDISSRV")
	if addr != "" {
		options.Addr = addr
	}
	user := os.Getenv("REDISUSER")
	if user != "" {
		options.Username = user
	}
	password := os.Getenv("REDISPASSWORD")
	if password != "" {
		options.Password = password
	}
	dbs := os.Getenv("REDISDB")
	if dbs != "" {
		if i, err := strconv.ParseInt(dbs, 10, 32); err == nil {
			options.DB = int(i)
		}
	}

	return options
}

// NewRdbClient create an instance of RdbClient
func NewRdbClient(options *redis.Options) (*RdbClient, error) {
	if options == nil {
		options = defaultOptions()
	}
	c := &RdbClient{
		client: redis.NewClient(options),
	}
	_, err := c.client.Ping(context.Background()).Result()
	if err != nil {
		return nil, ErrorServiceNotReachable
	}
	return c, nil
}

// Set set key value
func (c *RdbClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	s, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, s, expiration).Err()
}

// Get gets value by key
func (c *RdbClient) Get(ctx context.Context, key string, value interface{}) error {
	s, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrorKeyNotExists
		}
		return err
	}
	return json.Unmarshal([]byte(s), value)
}

// Expire set  key expire at given expiration duration
func (c *RdbClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	b, err := c.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		return err
	}
	if !b {
		return ErrorKeyNotExists
	}
	return nil
}

// ReadTimeout returns the ReadTimeout setting in redis client
func (c *RdbClient) ReadTimeout() time.Duration {
	return c.client.Options().ReadTimeout
}

// WriteTimeout returns the WriteTimeout setting in redis client
func (c *RdbClient) WriteTimeout() time.Duration {
	return c.client.Options().WriteTimeout
}
