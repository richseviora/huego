package pkg

import "context"

// BridgeRegistrationRequest represents the registration request payload
type BridgeRegistrationRequest struct {
	DeviceType        string `json:"devicetype"`
	GenerateClientKey bool   `json:"generateclientkey"`
}

// BridgeRegistrationResponse represents the registration response
type BridgeRegistrationResponse struct {
	Success map[string]string `json:"success"`
	Error   map[string]string `json:"error"`
}

// RegisterDevice sends a registration request to the Hue bridge
func (c *APIClient) RegisterDevice(ctx context.Context, appName, instanceName string) (*BridgeRegistrationResponse, error) {
	request := BridgeRegistrationRequest{
		DeviceType:        appName + "#" + instanceName,
		GenerateClientKey: true,
	}
	return Post[BridgeRegistrationResponse](ctx, "/api", request, c)
}
