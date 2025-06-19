package bridge

import (
	"encoding/json"
	"os"
	"time"
)

const bridgeCacheFile = "bridge_cache.json"

type BridgeCache struct {
	Bridges   []Bridge  `json:"bridges"`
	Timestamp time.Time `json:"timestamp"`
}

func LoadBridgeCache() (*BridgeCache, error) {
	data, err := os.ReadFile(bridgeCacheFile)
	if err != nil {
		return nil, err
	}
	var cache BridgeCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func SaveBridgeCache(cache *BridgeCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	return os.WriteFile(bridgeCacheFile, data, 0644)
}
