package pkg

import (
	"context"
	"errors"
	"github.com/richseviora/huego/internal/bridge"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/pkg/logger"
	client2 "github.com/richseviora/huego/pkg/resources/client"
)

func NewClientFromMDNS(logger logger.Logger) (client2.HueServiceClient, error) {
	bridges, err := bridge.DiscoverBridgesWithMDNS(logger)
	if err != nil {
		return nil, err
	}
	if len(bridges) == 0 {
		return nil, errors.New("no bridges found")
	}
	c := client.NewAPIClient(bridges[0].InternalIPAddress, client.EnvOnly, logger)
	err = c.Initialize(context.Background())
	if err != nil {
		return nil, err
	}
	return c, nil
}
