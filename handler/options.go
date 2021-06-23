package handler

import (
	"test.com/cache-server/cache"
	"test.com/cache-server/store"
)

// Option is func that set Handler option
type Option func(*Handler)

// WithCacher returns function sets the handler's cacher
func WithCacher(c cache.Cacher) Option {
	return func(h *Handler) {
		h.cache = c
	}
}

// WithStorer returns function sets the handler's storer
func WithStorer(s store.Storer) Option {
	return func(h *Handler) {
		h.store = s
	}
}
