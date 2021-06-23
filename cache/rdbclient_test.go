package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClient_Get(t *testing.T) {
	assert := assert.New(t)

	c, err := NewRdbClient(nil)
	if err != nil && errors.Is(err, ErrorServiceNotReachable) {
		t.Fatal("Can't create redis client, make sure server started")
	}

	key, value := uuid.New().String(), uuid.New().String()
	err = c.Set(context.Background(), key, value, 0)
	assert.NoError(err)

	var v string
	err = c.Get(context.Background(), key, &v)
	assert.NoError(err)
	assert.Equal(v, value)
}

func TestClient_Expire(t *testing.T) {
	assert := assert.New(t)

	c, err := NewRdbClient(nil)
	if err != nil && errors.Is(err, ErrorServiceNotReachable) {
		t.Fatal("Can't create redis client, make sure server started")
	}

	key, value := uuid.New().String(), uuid.New().String()
	err = c.Set(context.Background(), key, value, 0)
	assert.NoError(err)

	err = c.Expire(context.Background(), key, 0)
	assert.NoError(err)

	// the key expires
	var v string
	err = c.Get(context.Background(), key, &v)
	assert.Error(err)
	assert.Equal(v, "")
}
