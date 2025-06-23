package bridge

import (
	"encoding/json"
	"github.com/richseviora/huego/pkg/logger"
	"os"
	"time"
)

const bridgeCacheFile = "bridge_cache.json"

type BridgeCache struct {
	Bridges   map[string]Bridge `json:"bridges"`
	Timestamp time.Time         `json:"timestamp"`
	// Maps bridge IDs to application keys.
	ApplicationKeys map[string]string `json:"application_keys"`
}

type CacheManager struct {
	cache        *BridgeCache
	fileLocation string
	logger       logger.Logger
}

func NewCacheManager(fileLocation string, l logger.Logger) *CacheManager {
	return &CacheManager{
		fileLocation: fileLocation,
		logger:       l,
	}
}

func (c *CacheManager) Load() error {
	c.logger.Trace("loading cache from file", map[string]interface{}{
		"file": c.fileLocation,
	})
	data, err := os.ReadFile(bridgeCacheFile)
	if err != nil {
		c.logger.Error("failed to read cache file", map[string]interface{}{
			"error": err,
		})
		return err
	}
	if err := json.Unmarshal(data, c.cache); err != nil {
		c.logger.Error("failed to unmarshal JSON", map[string]interface{}{
			"error": err,
		})
		return err
	}
	c.logger.Trace("cache loaded", map[string]interface{}{})
	return nil
}

func (c *CacheManager) Save() error {
	c.logger.Trace("saving cache to file", map[string]interface{}{
		"file": c.fileLocation,
	})
	data, err := json.Marshal(c.cache)
	if err != nil {
		return err
	}
	return os.WriteFile(bridgeCacheFile, data, 0644)
}

func (c *CacheManager) GetBridgeAndKey(id string) (Bridge, string, error) {
	var bridge Bridge
	found := false
	for _, b := range c.cache.Bridges {
		if b.ID == id {
			bridge = b
			found = true
			break
		}
	}
	if !found {
		return Bridge{}, "", NoBridgeFoundError
	}
	result, ok := c.cache.ApplicationKeys[id]
	if !ok {
		return Bridge{}, "", NoKeyForBridgeError
	}
	return bridge, result, nil
}

func (c *CacheManager) UpdateBridgeData() error {
	newBridges, err := DiscoverBridges(c.logger)
	if err != nil {
		return err
	}

	if c.cache == nil {
		c.cache = &BridgeCache{}
	}
	if c.cache.Bridges == nil {
		c.cache.Bridges = make(map[string]Bridge)
	}
	if c.cache.ApplicationKeys == nil {
		c.cache.ApplicationKeys = make(map[string]string)
	}

	for _, newBridge := range newBridges {
		c.cache.Bridges[newBridge.ID] = newBridge
	}

	c.cache.Timestamp = time.Now()
	return c.Save()
}

func (c *CacheManager) SaveBridgeKeyForID(key, bridgeId string) error {
	c.cache.ApplicationKeys[bridgeId] = key
	return c.Save()
}

func (c *CacheManager) FindUnauthenticatedBridge() (Bridge, error) {
	for id := range c.cache.Bridges {
		_, exists := c.cache.ApplicationKeys[id]
		if !exists {
			return c.cache.Bridges[id], nil
		}
	}
	return Bridge{}, NoUnauthenticatedBridgeFoundError
}
