package main

import (
	"context"
	"fmt"
	"github.com/richseviora/huego/pkg/store"
	"time"

	"github.com/richseviora/huego/pkg"
	"github.com/richseviora/huego/pkg/resources"
)

func TestConnection(ipAddress string) error {
	client := pkg.NewAPIClient(ipAddress)
	err := client.Initialize(context.Background())
	if err != nil {
		fmt.Printf("Initialization Failed: %v\n", err)
		return err
	}
	lightService := resources.NewLightService(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lights, err := lightService.GetAllLights(ctx)
	if err != nil {
		fmt.Printf("Connection test failed: %v\n", err)
		return err
	}

	fmt.Printf("Connection successful! Found %d lights\n", len(lights.Data))
	for _, light := range lights.Data {
		fmt.Printf("- Light: %s (ID: %s)\n", light.Name, light.ID)
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
