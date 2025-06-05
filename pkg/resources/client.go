package resources

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/store"
	"net/http"
	"os"
	"time"
)

// InitMode defines the initialization mode for the API client
type InitMode int

const (
	// EnvOnly uses only environment variables for initialization
	EnvOnly InitMode = iota
	// EnvThenLocal tries environment variables first, then local storage
	EnvThenLocal
	// LocalOnly uses only local storage for initialization
	LocalOnly
)

// APIClient handles API communication
type APIClient struct {
	baseURL      string
	httpClient   *http.Client
	timeout      time.Duration
	keyStore     store.KeyStore
	initMode     InitMode
	LightService *LightService
	SceneService *SceneService
	RoomService  *RoomService
	ZoneService  *ZoneService
}

var (
	_ common.RequestProcessor = &APIClient{}
)

// ClientOption defines functional options for configuring the APIClient
type ClientOption func(*APIClient)

// NewAPIClient creates a new API client instance
func NewAPIClient(ipAddress string, initMode InitMode, opts ...ClientOption) *APIClient {
	var keyStore store.KeyStore = nil
	var err error = nil
	if initMode != EnvOnly {
		keyStore, err = store.NewDiskKeyStore("hue-keys.json")
		if err != nil {
			fmt.Printf("Failed to load key store: %v\n", err)
			panic(err)
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &APIClient{
		baseURL: "https://" + ipAddress,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
		timeout:  30 * time.Second,
		keyStore: keyStore,
		initMode: initMode,
	}
	client.SceneService = NewSceneService(client)
	client.LightService = light.NewLightService(client)
	client.RoomService = NewRoomService(client)
	client.ZoneService = NewZoneService(client)

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

func (c *APIClient) BaseURL() string {
	return c.baseURL
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

func (c *APIClient) getApplicationKey(ctx context.Context) (string, error) {
	switch c.initMode {
	case EnvOnly:
		key := os.Getenv("HUE_KEY")
		if key == "" {
			return "", store.ErrKeyNotFound
		}
		return key, nil
	case LocalOnly:
		key, err := c.keyStore.Get("application-key")
		if err != nil {
			return "", err
		}
		return key.(string), nil
	default: // EnvThenLocal
		if envKey := os.Getenv("HUE_KEY"); envKey != "" {
			return envKey, nil
		}
		key, err := c.keyStore.Get("application-key")
		if err != nil {
			return "", err
		}
		return key.(string), nil
	}
}

func (c *APIClient) setApplicationKey(ctx context.Context, key string) error {
	if c.initMode == EnvOnly {
		return errors.New("key store is disabled")
	}
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

// RegisterDevice sends a registration request to the Hue bridge
func (c *APIClient) RegisterDevice(ctx context.Context, appName, instanceName string) (*BridgeRegistrationResponseBody, error) {
	request := BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return common.Post[BridgeRegistrationResponseBody](ctx, "/api", request, c)
}
