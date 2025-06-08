package motion

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	common2 "github.com/richseviora/huego/internal/services/common"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/motion"
)

const basePath = "/clip/v2/resource/motion"

type Manager struct {
	client common.RequestProcessor
}

func (m *Manager) CollectionPath() string {
	return basePath
}

func (m *Manager) ResourcePath(id string) string {
	return basePath + "/" + id
}

func (m *Manager) GetAllMotion(ctx context.Context) (*common.ResourceList[motion.Data], error) {
	return handlers.Get[common.ResourceList[motion.Data]](ctx, m.CollectionPath(), m.client)
}

func (m *Manager) GetMotion(ctx context.Context, id string) (*motion.Data, error) {
	return handlers.GetSingularResource[motion.Data](id, m.ResourcePath(id), ctx, m.client, "motion")
}

var (
	_ motion.Service           = &Manager{}
	_ common2.ResourcePathable = &Manager{}
)

func NewManager(client common.RequestProcessor) *Manager {
	return &Manager{
		client: client,
	}
}
