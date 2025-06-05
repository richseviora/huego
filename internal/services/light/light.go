package light

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/light"
)

// LightService handles light-related API operations
type LightService struct {
	client common.RequestProcessor
}

var (
	_ light.LightManager = &LightService{}
)

// NewLightService creates a new LightService instance
func NewLightService(client common.RequestProcessor) *LightService {
	return &LightService{
		client: client,
	}
}

// GetLight retrieves a single light by its ID
func (s *LightService) GetLight(ctx context.Context, id string) (*light.Light, error) {
	result, err := common.Get[common.ResourceList[light.Light]](
		ctx,
		fmt.Sprintf("/clip/v2/resource/light/%s", id), s.client,
	)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 {
		return nil, fmt.Errorf("light not found")
	}
	room, err := common.FirstOrError(result)
	if err != nil {
		return nil, fmt.Errorf("light not found")
	}
	if room.ID != id {
		return nil, fmt.Errorf("light not found")
	}
	return room, nil
}

// GetAllLights retrieves all available lights
func (s *LightService) GetAllLights(ctx context.Context) (*common.ResourceList[light.Light], error) {
	return common.Get[common.ResourceList[light.Light]](ctx, "/clip/v2/resource/light", s.client)
}

func (s *LightService) UpdateLight(ctx context.Context, update light.LightUpdate) error {
	result, err := common.Put[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/light/"+update.ID, update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}
