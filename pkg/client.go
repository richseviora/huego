package pkg

import (
	"context"
	"errors"
	"github.com/richseviora/huego/pkg/resources"
	"github.com/richseviora/huego/pkg/store"
)

func NewClientFromMDNS() (*resources.APIClient, error) {
	bridges, err := store.DiscoverBridgesWithMDNS()
	if err != nil {
		return nil, err
	}
	if len(bridges) == 0 {
		return nil, errors.New("no bridges found")
	}
	client := resources.NewAPIClient(bridges[0].InternalIPAddress)
	err = client.Initialize(context.Background())
	if err != nil {
		return nil, err
	}
	return client, nil
}
