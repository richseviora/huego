package pkg

import (
	"github.com/richseviora/huego/internal/bridge"
	"github.com/richseviora/huego/pkg/logger"
	client2 "github.com/richseviora/huego/pkg/resources/client"
)

func NewClientProvider(logger logger.Logger) (client2.ClientProvider, error) {
	provider, err := bridge.NewBuilderWithoutPath(logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func NewClientProviderWithPath(path string, logger logger.Logger) (client2.ClientProvider, error) {
	provider, err := bridge.NewBuilderWithPath(path, logger)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
