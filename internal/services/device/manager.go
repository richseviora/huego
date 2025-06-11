package device

import (
	"context"
	"github.com/richseviora/huego/pkg/logger"

	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/device"
)

type DeviceManager struct {
	client common.RequestProcessor
	logger logger.Logger
}

var (
	_ device.Service = &DeviceManager{}
)

func NewDeviceManager(client common.RequestProcessor, logger logger.Logger) *DeviceManager {
	return &DeviceManager{
		client: client,
		logger: logger,
	}
}

func (m *DeviceManager) GetAllDevices(ctx context.Context) (*common.ResourceList[device.Data], error) {
	return handlers.Get[common.ResourceList[device.Data]](ctx, "/clip/v2/resource/device", m.client)
}

func (m *DeviceManager) GetDevice(ctx context.Context, id string) (*device.Data, error) {
	return handlers.GetSingularResource[device.Data](id, "/clip/v2/resource/device/"+id, ctx, m.client, "device")
}
