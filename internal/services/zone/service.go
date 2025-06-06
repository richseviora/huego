package zone

import (
	"context"
	"github.com/richseviora/huego/internal/client/handlers"
	"github.com/richseviora/huego/pkg/resources/common"
	zone2 "github.com/richseviora/huego/pkg/resources/zone"
)

type ZoneManager struct {
	client common.RequestProcessor
}

func NewZoneService(client common.RequestProcessor) *ZoneManager {
	return &ZoneManager{
		client: client,
	}
}

func (s *ZoneManager) GetAllZones(ctx context.Context) (*zone2.ZoneResponse, error) {
	return handlers.Get[zone2.ZoneResponse](ctx, "/clip/v2/resource/zone", s.client)
}

func (s *ZoneManager) GetZone(ctx context.Context, id string) (*zone2.ZoneData, error) {
	return handlers.GetSingularResource[zone2.ZoneData](id, "/clip/v2/resource/zone/"+id, ctx, s.client, "zone")
}

func (s *ZoneManager) CreateZone(ctx context.Context, zone *zone2.ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error) {
	return handlers.Post[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/zone", zone, s.client)
}

func (s *ZoneManager) UpdateZone(ctx context.Context, id string, zone *zone2.ZoneCreateOrUpdate) (*common.ResourceUpdateResponse, error) {
	return handlers.Put[common.ResourceUpdateResponse](ctx, "/clip/v2/resource/zone/"+id, zone, s.client)
}

func (s *ZoneManager) DeleteZone(ctx context.Context, id string) error {
	return handlers.Delete(ctx, "/clip/v2/resource/zone/"+id, s.client)
}
