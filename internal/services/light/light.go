package light

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/light"
)

// LightManager handles light-related API operations
type LightManager struct {
	client common.RequestProcessor
}

var (
	_ light.LightService = &LightManager{}
)

// NewLightService creates a new LightManager instance
func NewLightService(client common.RequestProcessor) *LightManager {
	return &LightManager{
		client: client,
	}
}

// GetLight retrieves a single light by its ID
func (s *LightManager) GetLight(ctx context.Context, id string) (*light.Light, error) {
	result, err := handlers.Get[common.ResourceList[light.Light]](
		ctx,
		fmt.Sprintf("/clip/v2/resource/light/%s", id), s.client,
	)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 {
		return nil, fmt.Errorf("light not found")
	}
	room, err := handlers.FirstOrError(result)
	if err != nil {
		return nil, fmt.Errorf("light not found")
	}
	if room.ID != id {
		return nil, fmt.Errorf("light not found")
	}
	return room, nil
}

// GetAllLights retrieves all available lights
func (s *LightManager) GetAllLights(ctx context.Context) (*common.ResourceList[light.Light], error) {
	return handlers.Get[common.ResourceList[light.Light]](ctx, "/clip/v2/resource/light", s.client)
}

func (s *LightManager) UpdateLight(ctx context.Context, update light.LightUpdate) error {
	result, err := handlers.Put[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/light/"+update.ID, update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}
