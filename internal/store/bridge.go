package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/richseviora/huego/pkg/logger"
	"net/http"
	"time"
)

type Bridge struct {
	ID                string `json:"id"`
	InternalIPAddress string `json:"internalipaddress"`
	Port              int    `json:"port"`
}

const url = "https://discovery.meethue.com"
const discoveryTimeout = time.Second * 5

func DiscoverBridgesWithMDNS(l logger.Logger) ([]Bridge, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolver: %v", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), discoveryTimeout)
	defer cancel()

	go func() {
		err = resolver.Browse(ctx, "_hue._tcp", "local.", entries)
		if err != nil {
			l.Error("Failed to browse for bridges", map[string]interface{}{
				"error": err,
			})
		}
	}()

	var bridges []Bridge
	for entry := range entries {
		port := entry.Port
		for _, ip := range entry.AddrIPv4 {
			bridges = append(bridges, Bridge{
				ID:                entry.Instance,
				InternalIPAddress: ip.String(),
				Port:              port,
			})
		}
	}

	return bridges, nil
}

func DiscoverBridges(logger logger.Logger) ([]Bridge, error) {
	// Try to load from cache first
	if cache, err := LoadBridgeCache(); err == nil {
		// Use cache if it's less than 1 hour old
		if time.Since(cache.Timestamp) < time.Hour {
			logger.Info("Found Bridges in Cache", map[string]interface{}{
				"bridges": len(cache.Bridges),
			})
			return cache.Bridges, nil
		}
	}

	// Try mDNS discovery first
	if bridges, err := DiscoverBridgesWithMDNS(logger); err == nil && len(bridges) > 0 {
		cache := &BridgeCache{
			Bridges:   bridges,
			Timestamp: time.Now(),
		}
		_ = SaveBridgeCache(cache)
		logger.Info("Found Bridges via MDNS", map[string]interface{}{
			"bridges": len(cache.Bridges),
		})
		return bridges, nil
	}

	// Fallback to HTTP discovery if mDNS fails

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bridges []Bridge
	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, err
	}

	// Save to cache
	cache := &BridgeCache{
		Bridges:   bridges,
		Timestamp: time.Now(),
	}
	_ = SaveBridgeCache(cache)

	logger.Info("Found Bridges via HTTP", map[string]interface{}{
		"bridges": len(bridges),
	})

	return bridges, nil
}
