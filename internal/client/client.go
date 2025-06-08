package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	behavior_instance2 "github.com/richseviora/huego/internal/services/behavior_instance"
	behavior_script2 "github.com/richseviora/huego/internal/services/behavior_script"
	motion2 "github.com/richseviora/huego/internal/services/motion"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/behavior_script"
	"github.com/richseviora/huego/pkg/resources/motion"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/richseviora/huego/internal/client/handlers"
	device2 "github.com/richseviora/huego/internal/services/device"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/internal/services/room"
	"github.com/richseviora/huego/internal/services/scene"
	zigbee_connectivity2 "github.com/richseviora/huego/internal/services/zigbee_connectivity"
	"github.com/richseviora/huego/internal/services/zone"
	"github.com/richseviora/huego/internal/store"
	"github.com/richseviora/huego/pkg/resources"
	"github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/device"
	light2 "github.com/richseviora/huego/pkg/resources/light"
	room2 "github.com/richseviora/huego/pkg/resources/room"
	scene2 "github.com/richseviora/huego/pkg/resources/scene"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	zone2 "github.com/richseviora/huego/pkg/resources/zone"
	"golang.org/x/time/rate"
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
	baseURL                   string
	httpClient                *http.Client
	timeout                   time.Duration
	keyStore                  store.KeyStore
	initMode                  InitMode
	limiter                   *rate.Limiter
	lightService              light2.LightService
	sceneService              scene2.SceneService
	roomService               room2.RoomService
	zoneService               zone2.ZoneService
	deviceService             device.Service
	zigbeeConnectivityService zigbee_connectivity.Service
	motionService             motion.Service
	behaviorInstanceService   behavior_instance.Service
	behaviorScriptService     behavior_script.Service
}

func (c *APIClient) BehaviorInstanceService() behavior_instance.Service {
	return c.behaviorInstanceService
}

func (c *APIClient) BehaviorScriptService() behavior_script.Service {
	return c.behaviorScriptService
}

func (c *APIClient) MotionService() motion.Service {
	return c.motionService
}

func (c *APIClient) ZigbeeConnectivityService() zigbee_connectivity.Service {
	return c.zigbeeConnectivityService
}

func (c *APIClient) DeviceService() device.Service {
	return c.deviceService
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
	baseUrl := ipAddress
	if !strings.HasPrefix(baseUrl, "https://") && !strings.HasPrefix(baseUrl, "http://") {
		baseUrl = "https://" + baseUrl
	}
	c := &APIClient{
		baseURL: baseUrl,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
		timeout:  30 * time.Second,
		keyStore: keyStore,
		initMode: initMode,
		limiter:  rate.NewLimiter(rate.Every(time.Second/10), 1),
	}
	c.sceneService = scene.NewSceneService(c)
	c.lightService = light.NewLightService(c)
	c.roomService = room.NewRoomService(c)
	c.zoneService = zone.NewZoneService(c)
	c.deviceService = device2.NewDeviceManager(c)
	c.zigbeeConnectivityService = zigbee_connectivity2.NewManager(c)
	c.motionService = motion2.NewManager(c)
	c.behaviorInstanceService = behavior_instance2.NewManager(c)
	c.behaviorScriptService = behavior_script2.NewManager(c)

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
	return createApplicationKey(ctx, c)
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

	err = c.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, client.ErrNotFound
	}
	if response.StatusCode == 503 {
		return nil, client.ErrServiceUnavailable
	}
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

func createApplicationKey(ctx context.Context, c *APIClient) error {
	res, err := c.registerDevice(ctx, "huego", "1234567890")
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

// registerDevice sends a registration request to the Hue bridge
func (c *APIClient) registerDevice(ctx context.Context, appName, instanceName string) (*resources.BridgeRegistrationResponseBody, error) {
	request := resources.BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return handlers.Post[resources.BridgeRegistrationResponseBody](ctx, "/api", request, c)
}
