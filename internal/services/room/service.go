package room

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/pkg/resources/common"
	room2 "github.com/richseviora/huego/pkg/resources/room"
)

const roomBasePath = "/clip/v2/resource/room"

func (s *RoomService) CollectionPath() string {
	return roomBasePath
}

func (s *RoomService) ResourcePath(id string) string {
	return roomBasePath + "/" + id
}

var (
	_ common.ResourcePathable = &RoomService{}
	_ room2.RoomService       = &RoomService{}
)

type RoomService struct {
	client common.RequestProcessor
}

func NewRoomService(client common.RequestProcessor) *RoomService {
	return &RoomService{
		client: client,
	}
}

func (s *RoomService) GetAllRooms(ctx context.Context) (*common.ResourceList[room2.RoomData], error) {
	return client.Get[common.ResourceList[room2.RoomData]](ctx, s.CollectionPath(), s.client)
}

func (s *RoomService) GetRoom(ctx context.Context, id string) (*room2.RoomData, error) {
	path := s.ResourcePath(id)
	return client.GetSingularResource[room2.RoomData](id, path, ctx, s.client, "room")
}

func (s *RoomService) UpdateRoom(ctx context.Context, update room2.RoomUpdate) error {
	result, err := client.Put[common.ResourceUpdateResponse](ctx, s.ResourcePath(update.ID), update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id string) error {
	err := client.Delete(ctx, s.ResourcePath(id), s.client)
	return err
}

func (s *RoomService) CreateRoom(ctx context.Context, create room2.RoomCreate) (*common.Reference, error) {
	return client.CreateResource(s.CollectionPath(), ctx, create, s.client, "room")
}
