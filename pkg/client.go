package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/richseviora/huego/pkg/store"
	"io"
	"net/http"
	"time"
)

// APIClient handles API communication
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	keyStore   store.KeyStore
}

// ClientOption defines functional options for configuring the APIClient
type ClientOption func(*APIClient)

// NewAPIClient creates a new API client instance
func NewAPIClient(ipAddress string, opts ...ClientOption) *APIClient {
	keyStore, err := store.NewDiskKeyStore("hue-keys.json")
	if err != nil {
		fmt.Printf("Failed to load key store: %v\n", err)
		panic(err)
	}
	client := &APIClient{
		baseURL: "http://" + ipAddress,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout:  30 * time.Second,
		keyStore: keyStore,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *APIClient) Initialize(ctx context.Context) error {
	key, err := c.getApplicationKey(ctx)
	if err != nil && !errors.Is(err, store.ErrKeyNotFound) {
		return err
	}
	if key != "" {
		// Key Set do nothing
		return nil
	}
	// Get Attempt to Set Key
	return CreateApplicationKey(ctx, c)
}

// Do executes an HTTP request and returns the response
func (c *APIClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	key, err := c.getApplicationKey(ctx)
	if err != nil && !errors.Is(err, store.ErrKeyNotFound) {
		return nil, err
	}
	if key != "" {
		req.Header.Set("hue-application-key", key)
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

func (c *APIClient) getApplicationKey(ctx context.Context) (string, error) {
	key, err := c.keyStore.Get("application-key")
	if err != nil {
		return "", err
	}
	return key.(string), nil
}

func (c *APIClient) setApplicationKey(ctx context.Context, key string) error {
	return c.keyStore.Set("application-key", key)
}

type CreateUserRequest struct {
	devicetype        string `json:"devicetype"`
	generateClientKey bool   `json:"generateclientkey"`
}

func CreateApplicationKey(ctx context.Context, c *APIClient) error {
	res, err := c.RegisterDevice(ctx, "huego", "1234567890")
	if err != nil {
		fmt.Printf("Failed to register device: %v\n", err)
		return err
	}
	for _, response := range *res {
		if response.Success != nil {
			return c.setApplicationKey(ctx, response.Success.Username)
		}
		if response.Error != nil {
			fmt.Printf("Failed to register device %v", response.Error)
			return errors.New("failed to register device")
		}
	}

	return errors.New("failed to register device")
}
