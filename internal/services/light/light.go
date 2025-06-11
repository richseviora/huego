package light

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	common2 "github.com/richseviora/huego/internal/services/common"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/light"
)

// LightManager handles light-related API operations
type LightManager struct {
	client common.RequestProcessor
	logger logger.Logger
}

var (
	_ light.LightService       = &LightManager{}
	_ common2.ResourcePathable = &LightManager{}
)

// NewLightService creates a new LightManager instance
func NewLightService(client common.RequestProcessor, logger logger.Logger) *LightManager {
	return &LightManager{
		client: client,
		logger: logger,
	}
}

func (s *LightManager) ResourcePath(id string) string {
	return "/clip/v2/resource/light/" + id
}

func (s *LightManager) CollectionPath() string {
	return "/clip/v2/resource/light"
}

// GetLight retrieves a single light by its ID
func (s *LightManager) GetLight(ctx context.Context, id string) (*light.Light, error) {
	return handlers.GetSingularResource[light.Light](id, s.ResourcePath(id), ctx, s.client, "light")

}

// GetAllLights retrieves all available lights
func (s *LightManager) GetAllLights(ctx context.Context) (*common.ResourceList[light.Light], error) {
	return handlers.Get[common.ResourceList[light.Light]](ctx, s.CollectionPath(), s.client)
}

func (s *LightManager) UpdateLight(ctx context.Context, update light.LightUpdate) error {
	_, err := handlers.UpdateResource[light.LightUpdate](s.ResourcePath(update.ID), ctx, update, s.client, "light")
	return err
}
