package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grandcat/zeroconf"
	"github.com/richseviora/huego/pkg/logger"
	"net/http"
	"time"
)

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
	var bridges []Bridge
	// Try mDNS discovery first
	if bridges, err := DiscoverBridgesWithMDNS(logger); err == nil && len(bridges) > 0 {
		logger.Info("Found Bridges via MDNS", map[string]interface{}{
			"bridges": len(bridges),
		})
		return bridges, nil
	}

	// Fallback to HTTP discovery if mDNS fails

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("failed to load bridge via external discovery", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&bridges); err != nil {
		return nil, err
	}

	logger.Info("Found Bridges via HTTP", map[string]interface{}{
		"bridges": len(bridges),
	})

	return bridges, nil
}
