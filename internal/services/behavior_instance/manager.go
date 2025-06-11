package behavior_instance

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	common2 "github.com/richseviora/huego/internal/services/common"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/common"
)

const basePath = "/clip/v2/resource/behavior_instance"

type Manager struct {
	client common.RequestProcessor
	logger logger.Logger
}

func (m *Manager) GetAllBehaviorInstances(ctx context.Context) (*common.ResourceList[behavior_instance.Data], error) {
	return handlers.Get[common.ResourceList[behavior_instance.Data]](ctx, m.CollectionPath(), m.client)
}

func (m *Manager) GetBehaviorInstance(ctx context.Context, id string) (*behavior_instance.Data, error) {
	return handlers.GetSingularResource[behavior_instance.Data](id, m.ResourcePath(id), ctx, m.client, "behavior_instance")
}

func (m *Manager) UpdateBehaviorInstance(ctx context.Context, id string, update behavior_instance.UpdateRequest) (*common.Reference, error) {
	return handlers.UpdateResource(m.ResourcePath(id), ctx, update, m.client, "behavior_instance")
}

func (m *Manager) CreateBehaviorInstance(ctx context.Context, create behavior_instance.CreateRequest) (*common.Reference, error) {
	return handlers.CreateResource[behavior_instance.CreateRequest](m.CollectionPath(), ctx, create, m.client, "behavior_instance")
}

func (m *Manager) DeleteBehaviorInstance(ctx context.Context, id string) error {
	return handlers.Delete(ctx, m.ResourcePath(id), m.client)
}

func (m *Manager) CollectionPath() string {
	return basePath
}

func (m *Manager) ResourcePath(id string) string {
	return basePath + "/" + id
}

var (
	_ behavior_instance.Service = &Manager{}
	_ common2.ResourcePathable  = &Manager{}
)

func NewManager(client common.RequestProcessor, logger logger.Logger) *Manager {
	return &Manager{
		client: client,
		logger: logger,
	}
}
