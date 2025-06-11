package room

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	common2 "github.com/richseviora/huego/internal/services/common"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/common"
	"github.com/richseviora/huego/pkg/resources/room"
)

const roomBasePath = "/clip/v2/resource/room"

func (s *RoomService) CollectionPath() string {
	return roomBasePath
}

func (s *RoomService) ResourcePath(id string) string {
	return roomBasePath + "/" + id
}

var (
	_ common2.ResourcePathable = &RoomService{}
	_ room.RoomService         = &RoomService{}
)

type RoomService struct {
	client common.RequestProcessor
	logger logger.Logger
}

func NewRoomService(client common.RequestProcessor, logger logger.Logger) *RoomService {
	return &RoomService{
		client: client,
		logger: logger,
	}
}

func (s *RoomService) GetAllRooms(ctx context.Context) (*common.ResourceList[room.RoomData], error) {
	return handlers.Get[common.ResourceList[room.RoomData]](ctx, s.CollectionPath(), s.client)
}

func (s *RoomService) GetRoom(ctx context.Context, id string) (*room.RoomData, error) {
	path := s.ResourcePath(id)
	return handlers.GetSingularResource[room.RoomData](id, path, ctx, s.client, "room")
}

func (s *RoomService) UpdateRoom(ctx context.Context, update room.RoomUpdate) error {
	_, err := handlers.UpdateResource[room.RoomUpdate](s.ResourcePath(update.ID), ctx, update, s.client, "room")
	return err
}

func (s *RoomService) DeleteRoom(ctx context.Context, id string) error {
	err := handlers.Delete(ctx, s.ResourcePath(id), s.client)
	return err
}

func (s *RoomService) CreateRoom(ctx context.Context, create room.RoomCreate) (*common.Reference, error) {
	return handlers.CreateResource(s.CollectionPath(), ctx, create, s.client, "room")
}
