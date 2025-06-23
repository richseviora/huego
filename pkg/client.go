package pkg

import (
	"github.com/richseviora/huego/internal/bridge"
	"github.com/richseviora/huego/pkg/logger"
	client2 "github.com/richseviora/huego/pkg/resources/client"
)

// NewClientWithoutPath constructs a HueServiceClient with the IP address and key supplied.
func NewClientWithoutPath(address, key string, logger logger.Logger) (client2.HueServiceClient, error) {
	if logger == nil {
		logger = NoOpLogger
	}
	provider, err := bridge.NewBuilderWithoutPath(logger)
	if err != nil {
		return nil, err
	}
	return provider.NewClientWithAddressAndKey(address, key)
}

// NewClientProviderWithPath returns a new Client Provider that you can use to generate an authenticated client.
func NewClientProviderWithPath(path string, logger logger.Logger) (client2.PersistentClientProvider, error) {
	if logger == nil {
		logger = NoOpLogger
	}
	provider, err := bridge.NewBuilderWithPath(path, logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

var NoOpLogger = logger.NoopLogger{}
