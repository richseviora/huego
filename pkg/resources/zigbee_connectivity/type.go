package zigbee_connectivity

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
)

type Channel struct {
	Status string `json:"status"`
	Value  string `json:"value"`
}
type Data struct {
	ID            string           `json:"id"`
	IDV1          string           `json:"id_v1,omitempty"`
	Owner         common.Reference `json:"owner"`
	Status        string           `json:"status"`
	MacAddress    string           `json:"mac_address"`
	Type          string           `json:"type"`
	Channel       Channel          `json:"channel,omitempty"`
	ExtendedPanID string           `json:"extended_pan_id,omitempty"`
}

var _ common.Identable = &Data{}

func (d Data) Identity() string {
	return d.ID
}

type Service interface {
	GetAllZigbeeConnectivity(ctx context.Context) (*common.ResourceList[Data], error)
	GetZigbeeConnectivity(ctx context.Context, id string) (*Data, error)
}
