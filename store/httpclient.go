package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	apiKey = `752f5fd3530c4c03825a225081c4b843`
)

// HttpClient defines the
type HttpClient struct {
	client *http.Client
}

// HttpClientOption is func that sets HttpClient option
type HttpClientOption func(*HttpClient) error

// NewHttpClient ...
func NewHttpClient(options ...HttpClientOption) (*HttpClient, error) {
	client := &HttpClient{
		client: http.DefaultClient,
	}
	for _, opt := range options {
		if err := opt(client); err != nil {
			log.WithError(err).Error("error configuring http client")
			return nil, err
		}
	}

	return client, nil
}

// Get run http get method
func (c *HttpClient) Get(ctx context.Context, url string, value interface{}) *ApiResponse {
	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
	}

	return c.do(req, value)
}

// Post run http post method
func (c *HttpClient) Post(ctx context.Context, url string, body []byte, value interface{}) *ApiResponse {
	req, err := c.newRequest(ctx, "POST", url, body)
	if err != nil {
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
	}

	return c.do(req, value)
}

func (c *HttpClient) newRequest(ctx context.Context, method, url string, data []byte) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		log.WithFields(log.Fields{"method": method, "url": url, "data": data}).WithError(err).Error("http new request failed")
		return nil, err
	}

	req.Header.Add("autopilotapikey", apiKey)
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *HttpClient) do(req *http.Request, value interface{}) *ApiResponse {
	resp, err := c.client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"req": req}).WithError(err).Error("http failed")
		return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("error reading body")
		}
		var apiResp ApiResponse
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "error unmarshal response"}
		}
		apiResp.StatusCode = resp.StatusCode
		return &apiResp
	}

	if value != nil {
		if resp.Header.Get("Content-Type") == "application/json" {
			err = json.NewDecoder(resp.Body).Decode(value)
		} else {
			err = fmt.Errorf("bad response header")
		}
		if err != nil {
			return &ApiResponse{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error()}
		}
	}

	return &ApiResponse{http.StatusOK, "OK", ""}
}
