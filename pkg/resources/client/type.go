package client

import (
	"github.com/richseviora/huego/pkg/resources/device"
	"github.com/richseviora/huego/pkg/resources/light"
	"github.com/richseviora/huego/pkg/resources/room"
	"github.com/richseviora/huego/pkg/resources/scene"
	"github.com/richseviora/huego/pkg/resources/zigbee_connectivity"
	"github.com/richseviora/huego/pkg/resources/zone"
)

type HueServiceClient interface {
	ZoneService() zone.ZoneService
	RoomService() room.RoomService
	SceneService() scene.SceneService
	LightService() light.LightService
	DeviceService() device.Service
	ZigbeeConnectivityService() zigbee_connectivity.Service
}
