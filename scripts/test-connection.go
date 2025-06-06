package main

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/internal/client"
	"github.com/richseviora/huego/internal/services/light"
	"github.com/richseviora/huego/internal/store"
	"time"
)

func TestConnection(ipAddress string) error {
	client := client.NewAPIClient(ipAddress, client.EnvThenLocal)
	err := client.Initialize(context.Background())
	if err != nil {
		fmt.Printf("Initialization Failed: %v\n", err)
		return err
	}
	lightService := light.NewLightService(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lights, err := lightService.GetAllLights(ctx)
	if err != nil {
		fmt.Printf("Connection test failed: %v\n", err)
		return err
	}

	fmt.Printf("Connection successful! Found %d lights\n", len(lights.Data))
	for _, light := range lights.Data {
		fmt.Printf("- Light: %s (ID: %s)\n", light.Metadata.Name, light.ID)
	}
	return nil
}

func main() {
	bridges, err := store.DiscoverBridges()
	if err != nil {
		fmt.Printf("Failed to discover bridges: %v\n", err)
		return
	}

	for _, bridge := range bridges {
		fmt.Printf("Found bridge: %s at %s\n", bridge.ID, bridge.InternalIPAddress)
	}

	err = TestConnection(bridges[0].InternalIPAddress)
	if err != nil {
		fmt.Printf("ABEND Connection test failed: %v\n", err)
	}
}
