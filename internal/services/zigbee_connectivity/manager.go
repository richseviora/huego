package zigbee_connectivity

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
)

type Manager struct {
	client common.RequestProcessor
}

var (
	_ zigbee_connectivity.Service = &Manager{}
)

func NewZigbeeConnectivityService(client common.RequestProcessor) *Manager {
	return &Manager{
		client: client,
	}
}

func (m Manager) GetAllZigbeeConnectivity(ctx context.Context) (*common.ResourceList[zigbee_connectivity.Data], error) {
	return handlers.Get[common.ResourceList[zigbee_connectivity.Data]](ctx, "/clip/v2/resource/zigbee_connectivity", m.client)
}

func (m Manager) GetZigbeeConnectivity(ctx context.Context, id string) (*zigbee_connectivity.Data, error) {
	return handlers.GetSingularResource[zigbee_connectivity.Data](id, "/clip/v2/resource/zigbee_connectivity/"+id, ctx, m.client, "zigbee_connectivity")
}
