package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grandcat/zeroconf"
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

func DiscoverBridgesWithMDNS() ([]Bridge, error) {
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
			fmt.Printf("Failed to browse: %v\n", err)
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

func DiscoverBridges() ([]Bridge, error) {
	// Try to load from cache first
	if cache, err := LoadBridgeCache(); err == nil {
		// Use cache if it's less than 1 hour old
		if time.Since(cache.Timestamp) < time.Hour {
			fmt.Printf("Found %d bridges in cache\n", len(cache.Bridges))
			return cache.Bridges, nil
		}
	}

	// Try mDNS discovery first
	if bridges, err := DiscoverBridgesWithMDNS(); err == nil && len(bridges) > 0 {
		cache := &BridgeCache{
			Bridges:   bridges,
			Timestamp: time.Now(),
		}
		_ = SaveBridgeCache(cache)
		fmt.Printf("Found %d bridges via mDNS\n", len(bridges))
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

	fmt.Printf("Found %d bridges via HTTP\n", len(bridges))
	return bridges, nil
}
