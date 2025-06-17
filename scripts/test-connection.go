package main

import (
	"context"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/internal/store"
	"github.com/richseviora/huego/pkg/logger"
	"time"
)

func TestConnection(ipAddress string, l logger.Logger) error {
	c := client.NewAPIClient(ipAddress, client.EnvThenLocal, l)
	err := c.Initialize(context.Background())
	if err != nil {
		l.Error("Initialization Failed: %v\n", map[string]interface{}{
			"error": err,
		})
		return err
	}
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
	bridges, err := store.DiscoverBridges(l)
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
