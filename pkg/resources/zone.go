package resources

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

type ZoneService struct {
	client *APIClient
}

func NewZoneService(client *APIClient) *ZoneService {
	return &ZoneService{
		client: client,
	}
}

func (s *ZoneService) GetAllZones(ctx context.Context) (*ZoneResponse, error) {
	return common.Get[ZoneResponse](ctx, "/clip/v2/resource/zone", s.client)
}

func (s *ZoneService) GetZone(ctx context.Context, id string) (*ZoneData, error) {
	return common.GetSingularResource[ZoneData](id, "/clip/v2/resource/zone/"+id, ctx, s.client, "zone")
}

func (s *ZoneService) CreateZone(ctx context.Context, zone *ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error) {
	return common.Post[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/zone", zone, s.client)
}

func (s *ZoneService) UpdateZone(ctx context.Context, id string, zone *ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error) {
	return common.Put[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/zone/"+id, zone, s.client)
}

func (s *ZoneService) DeleteZone(ctx context.Context, id string) error {
	return common.Delete(ctx, "/clip/v2/resource/zone/"+id, s.client)
}
