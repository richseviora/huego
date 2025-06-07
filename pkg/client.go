package pkg

import (
	"context"
	"errors"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/internal/store"
	client2 "github.com/richseviora/huego/pkg/resources/client"
)

func NewClientFromMDNS() (client2.HueServiceClient, error) {
	bridges, err := store.DiscoverBridgesWithMDNS()
	if err != nil {
		return nil, err
	}
	if len(bridges) == 0 {
		return nil, errors.New("no bridges found")
	}
	c := client.NewAPIClient(bridges[0].InternalIPAddress, client.EnvOnly)
	err = c.Initialize(context.Background())
	if err != nil {
		return nil, err
	}
	return c, nil
}
