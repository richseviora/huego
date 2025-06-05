package room

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
)

type RoomMetadata struct {
	Name      string      `json:"name"`
	Archetype common.Area `json:"archetype"`
}
type RoomData struct {
	ID       string             `json:"id"`
	IDV1     string             `json:"id_v1"`
	Children []common.Reference `json:"children"`
	Services []common.Reference `json:"services"`
	Metadata RoomMetadata       `json:"metadata"`
	Type     string             `json:"type"`
}

var (
	_ common.Identable = &RoomData{}
)

func (d RoomData) Identity() string {
	return d.ID
}

type RoomUpdate struct {
	ID       string              `json:"id"`
	Children *[]common.Reference `json:"children"`
	Metadata *RoomMetadata       `json:"metadata"`
}

type RoomCreate struct {
	Children []common.Reference `json:"children"`
	Metadata RoomMetadata       `json:"metadata"`
}

type RoomService interface {
	GetAllRooms(ctx context.Context) (*common.ResourceList[RoomData], error)
	GetRoom(ctx context.Context, id string) (*RoomData, error)
	UpdateRoom(ctx context.Context, update RoomUpdate) error
	DeleteRoom(ctx context.Context, id string) error
	CreateRoom(ctx context.Context, create RoomCreate) (*common.Reference, error)
}
