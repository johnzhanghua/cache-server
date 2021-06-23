package store

import (
	"context"
	"time"
)

// Storer interface defines the methods of
// Get/Upsert by key/value
type Storer interface {
	// Get method get response into the value
	Get(ctx context.Context, key string, value interface{}) *ApiResponse
	// Upsert method input data, and gets response into the value
	Upsert(ctx context.Context, key string, data, value interface{}) *ApiResponse

	Timeout() time.Duration
}
