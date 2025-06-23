package client

import (
	"github.com/richseviora/huego/pkg/resources/behavior_instance"
	"github.com/richseviora/huego/pkg/resources/behavior_script"
	"github.com/richseviora/huego/pkg/resources/device"
	"github.com/richseviora/huego/pkg/resources/light"
	"github.com/richseviora/huego/pkg/resources/motion"
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
	BehaviorInstanceService() behavior_instance.Service
	BehaviorScriptService() behavior_script.Service
	MotionService() motion.Service
}

type PersistentClientProvider interface {
	NewClientWithAddressAndKey(address string, key string) (HueServiceClient, error)
	NewClientWithNewBridge() (string, HueServiceClient, error)
	NewClientWithExistingBridge(bridgeId string) (HueServiceClient, error)
}

type ClientProvider interface {
	NewClientWithAddressAndKey(address string, key string) (HueServiceClient, error)
}
