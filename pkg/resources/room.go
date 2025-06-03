package resources

import (
	"context"
	"fmt"
)

type RoomMetadata struct {
	Name      string `json:"name"`
	Archetype Area   `json:"archetype"`
}
type RoomData struct {
	ID       string       `json:"id"`
	IDV1     string       `json:"id_v1"`
	Children []Reference  `json:"children"`
	Services []Reference  `json:"services"`
	Metadata RoomMetadata `json:"metadata"`
	Type     string       `json:"type"`
}

var (
	_ Identable = &RoomData{}
)

func (d RoomData) Identity() string {
	return d.ID
}

type RoomUpdate struct {
	ID       string        `json:"id"`
	Children *[]Reference  `json:"children"`
	Metadata *RoomMetadata `json:"metadata"`
}

type RoomCreate struct {
	Children []Reference  `json:"children"`
	Metadata RoomMetadata `json:"metadata"`
}

type RoomService struct {
	client *APIClient
}

const roomBasePath = "/clip/v2/resource/room"

func (s *RoomService) CollectionPath() string {
	return roomBasePath
}

func (s *RoomService) ResourcePath(id string) string {
	return roomBasePath + "/" + id
}

var (
	_ ResourcePathable = &RoomService{}
)

func NewRoomService(client *APIClient) *RoomService {
	return &RoomService{
		client: client,
	}
}

func (s *RoomService) GetAllRooms(ctx context.Context) (*ResourceList[RoomData], error) {
	return Get[ResourceList[RoomData]](ctx, s.CollectionPath(), s.client)
}

func (s *RoomService) GetRoom(ctx context.Context, id string) (*RoomData, error) {
	path := s.ResourcePath(id)
	return GetSingularResource[RoomData](id, path, ctx, s.client, "room")
}

func (s *RoomService) UpdateRoom(ctx context.Context, update RoomUpdate) error {
	result, err := Put[ResourceUpdateResponse](ctx, s.ResourcePath(update.ID), update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id string) error {
	err := Delete(ctx, s.ResourcePath(id), s.client)
	return err
}

func (s *RoomService) CreateRoom(ctx context.Context, create RoomCreate) (*Reference, error) {
	return CreateResource(s.CollectionPath(), ctx, create, s.client, "room")
}
