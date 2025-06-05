package zone

import (
	"context"
	"github.com/richseviora/huego/pkg/resources/common"
)

type ZoneResponse = common.ResourceList[ZoneData]

type ZoneMetadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
}

type ZoneData struct {
	ID       string             `json:"id"`
	IDV1     string             `json:"id_v1"`
	Children []common.Reference `json:"children"`
	Services []common.Reference `json:"services"`
	Metadata ZoneMetadata       `json:"metadata"`
	Type     string             `json:"type"`
}

func (z ZoneData) Identity() string {
	return z.ID
}

var _ common.Identable = &ZoneData{}

type ZoneCreateOrUpdate struct {
	Children []common.Reference `json:"children"`
	Metadata ZoneMetadata       `json:"metadata"`
}

type ZoneService interface {
	GetAllZones(ctx context.Context) (*ZoneResponse, error)
	GetZone(ctx context.Context, id string) (*ZoneData, error)
	CreateZone(ctx context.Context, zone *ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error)
	UpdateZone(ctx context.Context, id string, zone *ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error)
	DeleteZone(ctx context.Context, id string) error
}
