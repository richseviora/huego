package pkg

import "context"

// BridgeRegistrationRequest represents the registration request payload
type BridgeRegistrationRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

type BridgeRegistrationSuccess struct {
	Username  string `json:"username"`
	ClientKey string `json:"clientkey"`
}

type BridgeRegistrationError struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// BridgeRegistrationResponse represents the registration response
type BridgeRegistrationResponse struct {
	Success *BridgeRegistrationSuccess `json:"success"`
	Error   *BridgeRegistrationError   `json:"error"`
}

type BridgeRegistrationResponseBody = []BridgeRegistrationResponse

// RegisterDevice sends a registration request to the Hue bridge
func (c *APIClient) RegisterDevice(ctx context.Context, appName, instanceName string) (*BridgeRegistrationResponseBody, error) {
	request := BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return Post[BridgeRegistrationResponseBody](ctx, "/api", request, c)
}
