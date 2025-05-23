package resources

import (
	"context"
	"fmt"
)

type Area int

const (
	LivingRoom Area = iota
	Kitchen
	Dining
	Bedroom
	KidsBedroom
	Bathroom
	Nursery
	Recreation
	Office
	Gym
	Hallway
	Toilet
	FrontDoor
	Garage
	Terrace
	Garden
	Driveway
	Carport
	Home
	Downstairs
	Upstairs
	TopFloor
	Attic
	GuestRoom
	Staircase
	Lounge
	ManCave
	Computer
	Studio
	Music
	TV
	Reading
	Closet
	Storage
	LaundryRoom
	Balcony
	Porch
	Barbecue
	Pool
	Other
	InvalidArea
)

var AreaNames = [...]string{
	"living_room",
	"kitchen",
	"dining",
	"bedroom",
	"kids_bedroom",
	"bathroom",
	"nursery",
	"recreation",
	"office",
	"gym",
	"hallway",
	"toilet",
	"front_door",
	"garage",
	"terrace",
	"garden",
	"driveway",
	"carport",
	"home",
	"downstairs",
	"upstairs",
	"top_floor",
	"attic",
	"guest_room",
	"staircase",
	"lounge",
	"man_cave",
	"computer",
	"studio",
	"music",
	"tv",
	"reading",
	"closet",
	"storage",
	"laundry_room",
	"balcony",
	"porch",
	"barbecue",
	"pool",
	"other",
	"invalid_area",
}

// String returns the original snake-case token (or "Area(<n>)" if out of range).
func (a Area) String() string {
	if int(a) < 0 || int(a) >= len(AreaNames) {
		return fmt.Sprintf("Area(%d)", a)
	}
	return AreaNames[a]
}

// ParseArea converts a string like "living_room" back into the enum.
func ParseArea(s string) (Area, error) {
	for i, name := range AreaNames {
		if name == s {
			return Area(i), nil
		}
	}
	return InvalidArea, fmt.Errorf("unknown area: %q", s)
}

type Children struct {
	Rid   string `json:"rid"`
	Rtype string `json:"rtype"`
}
type RoomServices struct {
	Rid   string `json:"rid"`
	Rtype string `json:"rtype"`
}
type RoomMetadata struct {
	Name      string `json:"name"`
	Archetype Area   `json:"archetype"`
}
type RoomData struct {
	ID       string         `json:"id"`
	IDV1     string         `json:"id_v1"`
	Children []Children     `json:"children"`
	Services []RoomServices `json:"services"`
	Metadata RoomMetadata   `json:"metadata"`
	Type     string         `json:"type"`
}

type RoomUpdate struct {
	ID       string        `json:"id"`
	Children *[]Children   `json:"children"`
	Metadata *RoomMetadata `json:"metadata"`
}

type RoomCreate struct {
	Children []Children   `json:"children"`
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
	result, err := Get[ResourceList[RoomData]](ctx, fmt.Sprintf("/clip/v2/resource/room/%s", id), s.client)
	if err != nil {
		return nil, err
	}
	room, err := FirstOrError(result)
	if err != nil {
		return nil, fmt.Errorf("room not found")
	}
	if room.ID != id {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
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

func (s *RoomService) CreateRoom(ctx context.Context, create RoomCreate) (*Reference, error) {
	result, err := Post[ResourceUpdateResponse](ctx, "/clip/v2/resource/room", create, s.client)
	if err != nil {
		return nil, err
	}
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("failed to create room: %v", result.Errors)
	}
	return &result.Data[0], nil
}
