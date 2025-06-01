package resources

import "context"

type ZoneResponse = ResourceList[ZoneData]

type ZoneMetadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
}

type ZoneData struct {
	ID       string       `json:"id"`
	IDV1     string       `json:"id_v1"`
	Children []Reference  `json:"children"`
	Services []Reference  `json:"services"`
	Metadata ZoneMetadata `json:"metadata"`
	Type     string       `json:"type"`
}

func (z ZoneData) Identity() string {
	return z.ID
}

var _ Identable = &ZoneData{}

type ZoneCreateOrUpdate struct {
	Children []Reference  `json:"children"`
	Metadata ZoneMetadata `json:"metadata"`
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
	return Get[ZoneResponse](ctx, "/clip/v2/resource/zone", s.client)
}

func (s *ZoneService) GetZone(ctx context.Context, id string) (*ZoneData, error) {
	return GetSingularResource[ZoneData](id, "/clip/v2/resource/zone/"+id, ctx, s.client, "zone")
}

func (s *ZoneService) CreateZone(ctx context.Context, zone *ZoneData) (*ResourceUpdateResponse, error) {
	return Post[ResourceUpdateResponse](ctx, "/clip/v2/resource/zone", zone, s.client)
}

func (s *ZoneService) UpdateZone(ctx context.Context, zone *ZoneData) (*ResourceUpdateResponse, error) {
	return Put[ResourceUpdateResponse](ctx, "/clip/v2/resource/zone/"+zone.ID, zone, s.client)
}

func (s *ZoneService) DeleteZone(ctx context.Context, id string) error {
	return Delete(ctx, "/clip/v2/resource/zone/"+id, s.client)
}
