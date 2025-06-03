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

var _ Identable = &RoomData{}

type Identable interface {
	Identity() string
}

func (d RoomData) Identity() string {
	return d.ID
}

type RoomUpdate struct {
	ID       string        `json:"id"`
	Children *[]Child      `json:"children"`
	Metadata *RoomMetadata `json:"metadata"`
}

type RoomCreate struct {
	Children []Child      `json:"children"`
	Metadata RoomMetadata `json:"metadata"`
}

type RoomService struct {
	client *APIClient
}

func NewRoomService(client *APIClient) *RoomService {
	return &RoomService{
		client: client,
	}
}

func FirstOrError[T any](list *ResourceList[T]) (*T, error) {
	if list == nil || len(list.Data) == 0 {
		return nil, fmt.Errorf("resource not found")
	}
	return &list.Data[0], nil
}

func (s *RoomService) GetAllRooms(ctx context.Context) (*ResourceList[RoomData], error) {
	return Get[ResourceList[RoomData]](ctx, "/clip/v2/resource/room", s.client)
}

func (s *RoomService) GetRoom(ctx context.Context, id string) (*RoomData, error) {
	path := fmt.Sprintf("/clip/v2/resource/room/%s", id)
	return GetSingularResource[RoomData](id, path, ctx, s.client, "room")
}

func (s *RoomService) UpdateRoom(ctx context.Context, update RoomUpdate) error {
	result, err := Put[ResourceUpdateResponse](ctx, "/clip/v2/resource/room/"+update.ID, update, s.client)
	if err != nil {
		return err
	}
	if len(result.Errors) > 0 {
		return fmt.Errorf("failed to update light: %v", result.Errors)
	}
	return nil
}

func (s *RoomService) DeleteRoom(ctx context.Context, id string) error {
	err := Delete(ctx, "/clip/v2/resource/room/"+id, s.client)
	return err
}

func (s *RoomService) CreateRoom(ctx context.Context, create RoomCreate) (*Reference, error) {
	path := "/clip/v2/resource/room"
	return CreateResource(path, ctx, create, s.client, "room")
}

func CreateResource[T any](path string, ctx context.Context, create T, client *APIClient, resourceName string) (*Reference, error) {
	result, err := Post[ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to create resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}

func GetSingularResource[T Identable](id string, path string, ctx context.Context, client *APIClient, resourceName string) (*T, error) {
	result, err := Get[ResourceList[T]](ctx, path, client)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 {
		if result.Errors != nil && len(result.Errors) > 0 {
			return nil, fmt.Errorf(result.Errors[0].Description)
		}
		return nil, fmt.Errorf("resource ID %s of type %s not found", id, resourceName)
	}
	resource, err := FirstOrError[T](result)
	if err != nil {
		return nil, fmt.Errorf("resource ID %s of type %s not found", id, resourceName)
	}
	if (*resource).Identity() != id {
		return nil, fmt.Errorf("resource ID %s of type %s not matched", id, resourceName)
	}
	return resource, nil
}

func UpdateResource[T any](path string, ctx context.Context, create T, client *APIClient, resourceName string) (*Reference, error) {
	result, err := Put[ResourceUpdateResponse](ctx, path, create, client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to update resource %s: %v", resourceName, result.Errors)
	}
	return &result.Data[0], nil
}
