// Command pairing demonstrates node and device pairing workflows.
//
// It connects as an operator and exercises:
//   - Listing paired nodes and devices
//   - Node pairing: request, approve workflow
//   - Device pairing: list, approve/reject/remove
//   - Token rotation and revocation
//
// Usage:
//
//	go run ./examples/pairing
//
// Requires the mock server: go run ./examples/server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Pairing Example ===")
	fmt.Println()

	client := gateway.NewClient(
		gateway.WithToken("example-token"),
		gateway.WithRole(protocol.RoleOperator),
		gateway.WithScopes(
			protocol.ScopeOperatorRead,
			protocol.ScopeOperatorWrite,
			protocol.ScopeOperatorApprovals,
		),
		gateway.WithOnEvent(func(ev protocol.Event) {
			switch ev.EventName {
			case "node.pair.requested":
				fmt.Printf("  [event] Node pairing request: %s\n", string(ev.Payload))
			case "device.pair.requested":
				fmt.Printf("  [event] Device pairing request: %s\n", string(ev.Payload))
			}
		}),
	)
	defer client.Close()

	fmt.Println("Connecting...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}
	fmt.Println("Connected")
	fmt.Println()

	// --- Node Pairing ---
	fmt.Println("--- Node Pairing ---")

	fmt.Println("Listing paired nodes...")
	nodeList, err := client.NodePairList(ctx)
	if err != nil {
		fmt.Printf("NodePairList: %v\n", err)
	} else {
		fmt.Printf("Nodes: %s\n", formatJSON(nodeList))
	}

	fmt.Println("\nApproving node pairing request...")
	err = client.NodePairApprove(ctx, protocol.NodePairApproveParams{
		RequestID: "req-node-1",
	})
	if err != nil {
		fmt.Printf("NodePairApprove: %v\n", err)
	} else {
		fmt.Println("Node pairing approved")
	}

	// --- Device Pairing ---
	fmt.Println("\n--- Device Pairing ---")

	fmt.Println("Listing paired devices...")
	deviceList, err := client.DevicePairList(ctx)
	if err != nil {
		fmt.Printf("DevicePairList: %v\n", err)
	} else {
		fmt.Printf("Devices: %s\n", formatJSON(deviceList))
	}

	fmt.Println("\nApproving device pairing...")
	err = client.DevicePairApprove(ctx, protocol.DevicePairApproveParams{
		RequestID: "req-device-1",
	})
	if err != nil {
		fmt.Printf("DevicePairApprove: %v\n", err)
	} else {
		fmt.Println("Device pairing approved")
	}

	fmt.Println("\nRotating device token...")
	rotateResult, err := client.DeviceTokenRotate(ctx, protocol.DeviceTokenRotateParams{
		DeviceID: "device-1",
		Role:     "operator",
		Scopes:   []string{"operator.read", "operator.write"},
	})
	if err != nil {
		fmt.Printf("DeviceTokenRotate: %v\n", err)
	} else {
		fmt.Printf("Token rotated: %s\n", formatJSON(rotateResult))
	}

	fmt.Println("\nRemoving device...")
	err = client.DevicePairRemove(ctx, protocol.DevicePairRemoveParams{
		DeviceID: "device-1",
	})
	if err != nil {
		fmt.Printf("DevicePairRemove: %v\n", err)
	} else {
		fmt.Println("Device removed")
	}

	fmt.Println("\n=== Done ===")
}

func formatJSON(data json.RawMessage) string {
	var v any
	json.Unmarshal(data, &v)
	out, _ := json.MarshalIndent(v, "  ", "  ")
	return string(out)
}
