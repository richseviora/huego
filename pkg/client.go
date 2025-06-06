package pkg

import (
	"context"
	"errors"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/internal/store"
)

func NewClientFromMDNS() (*client.APIClient, error) {
	bridges, err := store.DiscoverBridgesWithMDNS()
	if err != nil {
		return nil, err
	}
	if len(bridges) == 0 {
		return nil, errors.New("no bridges found")
	}
	client := client.NewAPIClient(bridges[0].InternalIPAddress, client.EnvOnly)
	err = client.Initialize(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}
