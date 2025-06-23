package client

import (
	"context"
	"errors"
	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources"
	client2 "github.com/richseviora/huego/pkg/resources/client"
	"github.com/richseviora/huego/pkg/resources/common"
	"net/http"
)

type BridgeRegistrationClient struct {
	logger     logger.Logger
	baseURL    string
	httpClient *http.Client
}

func NewBridgeRegistrationClient(baseURL string, logger logger.Logger) *BridgeRegistrationClient {
	return &BridgeRegistrationClient{
		baseURL:    baseURL,
		logger:     logger,
		httpClient: NewHTTPClient(),
	}
}

func (c *BridgeRegistrationClient) Logger() logger.Logger {
	return c.logger
}

func (c *BridgeRegistrationClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, client2.ErrNotFound
	}
	if response.StatusCode == 503 {
		return nil, client2.ErrServiceUnavailable
	}
	if response.StatusCode == 403 {
		return nil, client2.ErrUnauthorized
	}
	return response, err
}

func (c *BridgeRegistrationClient) BaseURL() string {
	return c.baseURL
}

var (
	_ common.RequestProcessor = &BridgeRegistrationClient{}
)

func (c *BridgeRegistrationClient) RegisterDevice(ctx context.Context, appName, instanceName string) (string, error) {
	if appName == "" {
		appName = "huego"
	}
	if instanceName == "" {
		instanceName = "default"
	}
	response, err := c.registerDevice(ctx, appName, instanceName)
	if err != nil {
		return "", err
	}
	registrationError := (*response)[0].Error
	if registrationError != nil {
		if registrationError.Description == "link button not pressed" {
			return "", LinkButtonNotPressedError
		}
	}
	registration := (*response)[0].Success
	return registration.Username, nil
}

// registerDevice sends a registration request to the Hue bridge
func (c *BridgeRegistrationClient) registerDevice(ctx context.Context, appName, instanceName string) (*resources.BridgeRegistrationResponseBody, error) {
	request := resources.BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return handlers.Post[resources.BridgeRegistrationResponseBody](ctx, "/api", request, c)
}

var LinkButtonNotPressedError = errors.New("link button not pressed")
