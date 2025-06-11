package behavior_script

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	common2 "github.com/richseviora/huego/internal/services/common"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/behavior_script"
	"github.com/richseviora/huego/pkg/resources/common"
)

const basePath = "/clip/v2/resource/behavior_script"

type Manager struct {
	client common.RequestProcessor
	logger logger.Logger
}

func (m *Manager) GetAllBehaviorScripts(ctx context.Context) (*common.ResourceList[behavior_script.Data], error) {
	return handlers.Get[common.ResourceList[behavior_script.Data]](ctx, m.CollectionPath(), m.client)
}

func (m *Manager) GetBehaviorScript(ctx context.Context, id string) (*behavior_script.Data, error) {
	return handlers.GetSingularResource[behavior_script.Data](id, m.ResourcePath(id), ctx, m.client, "behavior_script")
}

func (m *Manager) CollectionPath() string {
	return basePath
}

func (m *Manager) ResourcePath(id string) string {
	return basePath + "/" + id
}

var (
	_ behavior_script.Service  = &Manager{}
	_ common2.ResourcePathable = &Manager{}
)

func NewManager(client common.RequestProcessor, logger logger.Logger) *Manager {
	return &Manager{
		client: client,
		logger: logger,
	}
}
