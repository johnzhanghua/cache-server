package store

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"test.com/cache-server/utils"
)

const (
	apiTimeout = 30 * time.Second
)

// AutoPilotStore implements Storer interface
type AutoPilotStore struct {
	client *HttpClient
}

// ApiResponse is the response message from autopilot api , and also send to user
type ApiResponse struct {
	StatusCode int    `json:"status,omitempty"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

// NewAutoPilotStore create instance of HttpStore,
// and return as a Storer interface
func NewAutoPilotStore(httpOptions ...HttpClientOption) (Storer, error) {
	c, err := NewHttpClient(httpOptions...)
	if err != nil {
		return nil, err
	}
	return &AutoPilotStore{
		client: c,
	}, nil
}

// Get gets value from store
func (s *AutoPilotStore) Get(ctx context.Context, key string, value interface{}) *ApiResponse {
	// from the key get the url
	url, err := utils.GetURLFromKey(key)
	if err != nil {
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
	}
	return s.client.Get(ctx, url, value)
}

// Upsert creates/updates and save key/value to store
func (s *AutoPilotStore) Upsert(ctx context.Context, key string, data, resp interface{}) *ApiResponse {
	url, err := utils.GetPostURLFromKey(key)
	if err != nil {
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
	}
	body, err := json.Marshal(data)
	if err != nil {
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "error marshalling json"}
	}
	return s.client.Post(ctx, url, body, resp)
}

// Timeout returns the store api timeout value
func (s *AutoPilotStore) Timeout() time.Duration {
	return apiTimeout
}
