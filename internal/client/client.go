package client

import (
	"context"
	"crypto/tls"
	"errors"
	behavior_instance2 "github.com/richseviora/huego/internal/services/behavior_instance"
	behavior_script2 "github.com/richseviora/huego/internal/services/behavior_script"
	motion2 "github.com/richseviora/huego/internal/services/motion"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/behavior_script"
	"github.com/richseviora/huego/pkg/resources/motion"
	"net/http"
	"strings"
	"time"

	device2 "github.com/richseviora/huego/internal/services/device"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/internal/services/room"
	"github.com/richseviora/huego/internal/services/scene"
	zigbee_connectivity2 "github.com/richseviora/huego/internal/services/zigbee_connectivity"
	"github.com/richseviora/huego/internal/services/zone"
	"github.com/richseviora/huego/internal/store"
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
	logger                    logger.Logger
	baseURL                   string
	applicationKey            string
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

func (c *APIClient) Logger() logger.Logger {
	return c.logger
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
func NewAPIClient(ipAddress string, applicationKey string, logger logger.Logger, opts ...ClientOption) *APIClient {
	baseUrl := ipAddress
	if !strings.HasPrefix(baseUrl, "https://") && !strings.HasPrefix(baseUrl, "http://") {
		baseUrl = "https://" + baseUrl
	}
	c := &APIClient{
		logger:         logger,
		baseURL:        baseUrl,
		httpClient:     NewHTTPClient(),
		timeout:        30 * time.Second,
		applicationKey: applicationKey,
		limiter:        rate.NewLimiter(rate.Every(time.Second/10), 1),
	}
	c.sceneService = scene.NewSceneService(c, c.logger)
	c.lightService = light.NewLightService(c, c.logger)
	c.roomService = room.NewRoomService(c, c.logger)
	c.zoneService = zone.NewZoneService(c, c.logger)
	c.deviceService = device2.NewDeviceManager(c, c.logger)
	c.zigbeeConnectivityService = zigbee_connectivity2.NewManager(c, c.logger)
	c.motionService = motion2.NewManager(c, c.logger)
	c.behaviorInstanceService = behavior_instance2.NewManager(c, c.logger)
	c.behaviorScriptService = behavior_script2.NewManager(c, c.logger)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}
}

func (c *APIClient) BaseURL() string {
	return c.baseURL
}

// Do executes an HTTP request and returns the response
func (c *APIClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	key := c.applicationKey
	req.Header.Set("hue-application-key", key)

	err := c.limiter.Wait(ctx)
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

func createApplicationKey(ctx context.Context, c *BridgeRegistrationClient) (string, error) {
	res, err := c.registerDevice(ctx, "huego", "1234567890")
	if err != nil {
		c.Logger().Error("Failed to register device", map[string]interface{}{
			"error": err,
		})
		return "", err
	}
	for _, response := range *res {
		if response.Success != nil {
			return response.Success.ClientKey, nil
		}
		if response.Error != nil {
			c.Logger().Error("Failed to register device", map[string]interface{}{
				"error": err,
			})
			return "", errors.New("failed to register device")
		}
	}

	return "", errors.New("failed to register device")
}
