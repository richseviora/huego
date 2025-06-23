package bridge

import (
	"context"
	"errors"
	client2 "github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/pkg/logger"
	"github.com/richseviora/huego/pkg/resources/client"
)

type Builder struct {
	FileLocation  string
	BridgeManager *CacheManager
	Logger        logger.Logger
}

func NewBuilderWithPath(fileLocation string, logger logger.Logger) (client.PersistentClientProvider, error) {
	res := &Builder{
		FileLocation:  fileLocation,
		BridgeManager: NewCacheManager(fileLocation, logger),
		Logger:        logger,
	}
	err := res.BridgeManager.Load()
	if err != nil {
		res.Logger.Error("Failed to load cache: %v\n", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}
	return res, nil
}

func NewBuilderWithoutPath(logger logger.Logger) (client.ClientProvider, error) {
	return &Builder{
		Logger: logger,
	}, nil
}

func (b *Builder) NewClientWithAddressAndKey(address string, key string) (client.HueServiceClient, error) {
	return client2.NewAPIClient(address, key, b.Logger), nil
}

func (b *Builder) NewClientWithNewBridge() (string, client.HueServiceClient, error) {
	if b.FileLocation != "" {
		return "", nil, NoFileLocationError
	}
	bridge, err := b.BridgeManager.FindUnauthenticatedBridge()
	if err != nil {
		return "", nil, err
	}
	r := client2.NewBridgeRegistrationClient(bridge.InternalIPAddress, b.Logger)
	key, err := r.RegisterDevice(context.Background(), "huego", "default")
	if err != nil {
		b.Logger.Error("Failed to register device", map[string]interface{}{
			"error":    err,
			"bridgeID": bridge.ID,
			"bridgeIP": bridge.InternalIPAddress,
		})
		return "", nil, err
	}
	b.Logger.Info("Registered device", map[string]interface{}{
		"bridgeID": bridge.ID,
	})
	err = b.BridgeManager.SaveBridgeKeyForID(key, bridge.ID)
	if err != nil {
		return "", nil, err
	}
	b.Logger.Trace("Saved key", map[string]interface{}{
		"bridgeID": bridge.ID,
	})
	return bridge.ID, client2.NewAPIClient(bridge.InternalIPAddress, key, b.Logger), nil
}

func (b *Builder) NewClientWithExistingBridge(bridgeId string) (client.HueServiceClient, error) {
	if b.FileLocation != "" {
		return nil, NoFileLocationError
	}
	bridge, key, err := b.BridgeManager.GetBridgeAndKey(bridgeId)
	if err != nil {
		return nil, err
	}
	return client2.NewAPIClient(bridge.InternalIPAddress, key, b.Logger), nil
}

var (
	_ client.PersistentClientProvider = &Builder{}
	_ client.ClientProvider           = &Builder{}
)

var NoFileLocationError = errors.New("no file location provided")
var NoWriteAccessError = errors.New("no write access to file location")
var NoKeyForBridgeError = errors.New("no key for bridge")
var NoBridgeFoundError = errors.New("no bridge found")
var NoUnauthenticatedBridgeFoundError = errors.New("no unauthenticated bridge found")
