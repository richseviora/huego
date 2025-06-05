package client

import (
	"github.com/richseviora/huego/pkg/resources/light"
	"github.com/richseviora/huego/pkg/resources/room"
	"github.com/richseviora/huego/pkg/resources/scene"
	"github.com/richseviora/huego/pkg/resources/zone"
)

type HueServiceClient interface {
	ZoneService() zone.ZoneService
	RoomService() room.RoomService
	SceneService() scene.SceneService
	LightService() light.LightService
}
