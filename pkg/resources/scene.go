package resources

import (
	"context"
	"fmt"
)

type Action struct {
	Target Reference `json:"target"`
	Action struct {
		On               LightOn          `json:"on"`
		Dimming          Dimming          `json:"dimming"`
		ColorTemperature ColorTemperature `json:"color_temperature"`
	}
}

type Palette struct{}

type Scene struct {
	ID       string   `json:"id"`
	IDV1     string   `json:"id_v1"`
	Actions  []Action `json:"actions"`
	Palette  Palette  `json:"palette"`
	Recall   struct{} `json:"recall"`
	Metadata struct {
		Name  string    `json:"name"`
		Image Reference `json:"image"`
	}
	// Can either be room or zone.
	Group       Reference `json:"group"`
	Speed       float64   `json:"speed"`
	AutoDynamic bool      `json:"auto_dynamic"`
	Status      struct {
		Active     string `json:"active"`
		LastRecall string `json:"last_recall"`
	}
	Type string `json:"type"`
}

type SceneService struct {
	client *APIClient
}

func NewSceneService(client *APIClient) *SceneService {
	return &SceneService{
		client: client,
	}
}

func (s *SceneService) GetAllScenes(ctx context.Context) (*ResourceList[Scene], error) {
	return Get[ResourceList[Scene]](ctx, "/clip/v2/resource/scene", s.client)
}

func (s *SceneService) GetScene(ctx context.Context, id string) (*Scene, error) {
	return Get[Scene](ctx, fmt.Sprintf("/clip/v2/resource/scene/%s", id), s.client)
}
