package main

import (
	"context"
	"github.com/richseviora/huego/internal/bridge"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/pkg/logger"
	"os"
	"time"
)

func TestConnection(ipAddress string, l logger.Logger) error {
	key := os.Getenv("HUE_KEY")
	if key == "" {
		l.Error("HUE_KEY not set")
		return nil
	}

	c := client.NewAPIClient(ipAddress, key, l)

	lightService := light.NewLightService(c, logger.NoopLogger{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lights, err := lightService.GetAllLights(ctx)
	if err != nil {
		l.Error("Connection Test Failed: %v\n", map[string]interface{}{
			"error": err,
		})
		return err
	}

	l.Info("Connection Test Successful", map[string]interface{}{
		"lights": len(lights.Data),
	})
	for _, light := range lights.Data {
		l.Info("Light", map[string]interface{}{"name": light.Metadata.Name, "id": light.ID})
	}
	return nil
}

func main() {
	l := logger.NewLogger()
	bridges, err := bridge.DiscoverBridges(l)
	if err != nil {
		l.Error("Failed to discover bridges: %v\n", map[string]interface{}{
			"error": err,
		})
		return
	}

	for _, bridge := range bridges {
		l.Info("Found bridge", map[string]interface{}{
			"bridge": bridge,
		})
	}

	err = TestConnection(bridges[0].InternalIPAddress, l)
	if err != nil {
		l.Error("ABEND Connection test failed", map[string]interface{}{
			"error": err,
		})
	}
}
