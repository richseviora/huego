package resources

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
