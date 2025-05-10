package resources

import (
	"context"
	"fmt"

	"github.com/richseviora/huego/pkg"
)

// Light represents the light resource data
type Light struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	State      bool   `json:"state"`
	Brightness int    `json:"brightness"`
}

type ResourceList[T any] struct {
	Data   []T      `json:"data"`
	Errors []string `json:"errors"`
}

// LightService handles light-related API operations
type LightService struct {
	client *pkg.APIClient
}

// NewLightService creates a new LightService instance
func NewLightService(client *pkg.APIClient) *LightService {
	return &LightService{
		client: client,
	}
}

// GetLight retrieves a single light by its ID
func (s *LightService) GetLight(ctx context.Context, id string) (*Light, error) {
	return pkg.Get[Light](ctx, fmt.Sprintf("/light/%s", id), s.client)
}

// GetAllLights retrieves all available lights
func (s *LightService) GetAllLights(ctx context.Context) (*ResourceList[Light], error) {
	return pkg.Get[ResourceList[Light]](ctx, "/light", s.client)
}
