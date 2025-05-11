package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/richseviora/huego/pkg/store"
	"io"
	"net/http"
	"time"
)

// APIClient handles API communication
type APIClient struct {
	baseURL        string
	httpClient     *http.Client
	timeout        time.Duration
	applicationKey string
	keyStore       store.KeyStore
}

// ClientOption defines functional options for configuring the APIClient
type ClientOption func(*APIClient)

// NewAPIClient creates a new API client instance
func NewAPIClient(ipAddress string, opts ...ClientOption) *APIClient {
	client := &APIClient{
		baseURL: "http://" + ipAddress,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout:  30 * time.Second,
		keyStore: store.NewInMemoryKeyStore(),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Do executes an HTTP request and returns the response
func (c *APIClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if c.applicationKey != "" {
		req.Header.Set("hue-application-key", c.applicationKey)
	}

	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

// Get performs a GET request and unmarshals the response into the provided type
func Get[T any](ctx context.Context, path string, c *APIClient) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result T
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		fmt.Printf("Failed to decode response: %v\n", err)
		fmt.Printf("Response body: %s\n", string(bodyBytes))
		return nil, err
	}

	return &result, nil
}

// Post performs a POST request and unmarshals the response into the provided type
func Post[T any](ctx context.Context, path string, body interface{}, c *APIClient) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
