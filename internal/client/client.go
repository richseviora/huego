package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/internal/services/room"
	"github.com/richseviora/huego/internal/services/scene"
	"github.com/richseviora/huego/internal/services/zone"
	"github.com/richseviora/huego/internal/store"
	"github.com/richseviora/huego/pkg/resources"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/common"
	light2 "github.com/richseviora/huego/pkg/resources/light"
	room2 "github.com/richseviora/huego/pkg/resources/room"
	scene2 "github.com/richseviora/huego/pkg/resources/scene"
	zone2 "github.com/richseviora/huego/pkg/resources/zone"

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
	lightService light2.LightService
	sceneService scene2.SceneService
	roomService  room2.RoomService
	zoneService  zone2.ZoneService
}

func (c *APIClient) ZoneService() zone2.ZoneService {
	return c.zoneService
}

func (c *APIClient) RoomService() room2.RoomService {
	return c.roomService
}

func (c *APIClient) SceneService() scene2.SceneService {
	return c.sceneService
}

func (c *APIClient) LightService() light2.LightService {
	return c.lightService
}

var (
	_ common.RequestProcessor = &APIClient{}
	_ client.HueServiceClient = &APIClient{}
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
	c := &APIClient{
		baseURL: "https://" + ipAddress,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
		timeout:  30 * time.Second,
		keyStore: keyStore,
		initMode: initMode,
	}
	c.sceneService = scene.NewSceneService(c)
	c.lightService = light.NewLightService(c)
	c.roomService = room.NewRoomService(c)
	c.zoneService = zone.NewZoneService(c)

	for _, opt := range opts {
		opt(c)
	}

	return c
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
	response, err := c.httpClient.Do(req)
	if response.StatusCode == 403 {
		return nil, client.ErrUnauthorized
	}
	return response, err
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
func (c *APIClient) RegisterDevice(ctx context.Context, appName, instanceName string) (*resources.BridgeRegistrationResponseBody, error) {
	request := resources.BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return handlers.Post[resources.BridgeRegistrationResponseBody](ctx, "/api", request, c)
}
