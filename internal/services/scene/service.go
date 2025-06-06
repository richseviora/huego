package scene

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/pkg/resources/common"
	scene2 "github.com/richseviora/huego/pkg/resources/scene"
)

var (
	_ scene2.SceneService = &SceneManager{}
)

type SceneManager struct {
	client common.RequestProcessor
}

func NewSceneService(client common.RequestProcessor) *SceneManager {
	return &SceneManager{
		client: client,
	}
}

func (s *SceneManager) GetAllScenes(ctx context.Context) (*common.ResourceList[scene2.SceneData], error) {
	return client.Get[common.ResourceList[scene2.SceneData]](ctx, "/clip/v2/resource/scene", s.client)
}

func (s *SceneManager) GetScene(ctx context.Context, id string) (*scene2.SceneData, error) {
	path := fmt.Sprintf("/clip/v2/resource/scene/%s", id)
	return client.GetSingularResource[scene2.SceneData](id, path, ctx, s.client, "scene")
}

func (s *SceneManager) UpdateScene(ctx context.Context, id string, scene scene2.SceneUpdate) (*common.Reference, error) {
	url := fmt.Sprintf("/clip/v2/resource/scene/%s", id)
	return client.UpdateResource(url, ctx, scene, s.client, "scene")
}

func (s *SceneManager) CreateScene(ctx context.Context, scene scene2.SceneCreate) (*common.Reference, error) {
	return client.CreateResource("/clip/v2/resource/scene", ctx, scene, s.client, "scene")
}

func (s *SceneManager) DeleteScene(ctx context.Context, id string) error {
	return client.Delete(ctx, fmt.Sprintf("/clip/v2/resource/scene/%s", id), s.client)
}
