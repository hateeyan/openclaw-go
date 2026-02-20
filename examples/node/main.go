// Command node demonstrates connecting as a node (capability host).
//
// Nodes provide execution capabilities to the gateway. This example:
//   - Connects as a node role
//   - Declares its capabilities and commands
//   - Handles invoke requests from the gateway
//   - Sends node events back to the gateway
//
// Usage:
//
//	go run ./examples/node
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

	fmt.Println("=== OpenClaw Node Example ===")
	fmt.Println()

	// Connect as a node with an invoke handler.
	client := gateway.NewClient(
		gateway.WithToken("node-token"),
		gateway.WithRole(protocol.RoleNode),
		gateway.WithScopes(protocol.ScopeOperatorRead),
		gateway.WithCaps("exec", "fs"),
		gateway.WithCommands("shell", "read_file", "write_file"),
		gateway.WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
			// Handle commands sent to this node by the gateway.
			fmt.Printf("  [invoke] Command: %s, ID: %s\n", inv.Command, inv.ID)
			fmt.Printf("  [invoke] Params: %s\n", string(inv.Params))

			// Return a success response.
			payload, _ := json.Marshal(map[string]any{
				"exitCode": 0,
				"output":   "command executed successfully",
			})
			return protocol.InvokeResponse{
				Type:    "invoke-res",
				ID:      inv.ID,
				OK:      true,
				Payload: payload,
			}
		}),
		gateway.WithOnEvent(func(ev protocol.Event) {
			fmt.Printf("  [event] %s\n", ev.EventName)
		}),
	)
	defer client.Close()

	fmt.Println("Connecting as node...")
	if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
		log.Fatalf("Connect: %v", err)
	}

	hello := client.Hello()
	fmt.Printf("Connected (protocol=%d)\n\n", hello.Protocol)

	// Request pairing with the gateway.
	fmt.Println("Requesting node pairing...")
	pairResult, err := client.NodePairRequest(ctx, protocol.NodePairRequestParams{
		NodeID:      "example-node-1",
		DisplayName: "Example Node",
		Platform:    "darwin",
		Version:     "1.0.0",
		Caps:        []string{"exec", "fs"},
		Commands:    []string{"shell", "read_file", "write_file"},
	})
	if err != nil {
		fmt.Printf("NodePairRequest: %v (may not be implemented in mock)\n", err)
	} else {
		fmt.Printf("Pair result: %s\n", string(pairResult))
	}

	// Send a node event.
	fmt.Println("\nSending node event...")
	err = client.NodeEvent(ctx, protocol.NodeEventParams{
		Event:       "node.ready",
		PayloadJSON: `{"caps":["exec","fs"],"status":"ready"}`,
	})
	if err != nil {
		fmt.Printf("NodeEvent: %v (may not be implemented in mock)\n", err)
	} else {
		fmt.Println("Node event sent")
	}

	// Send an invoke result (simulating a completed command).
	fmt.Println("\nSending invoke result...")
	resultPayload, _ := json.Marshal(map[string]any{
		"exitCode": 0,
		"output":   "hello world",
	})
	err = client.NodeInvokeResult(ctx, protocol.NodeInvokeResultParams{
		ID:      "invoke-1",
		NodeID:  "example-node-1",
		OK:      true,
		Payload: resultPayload,
	})
	if err != nil {
		fmt.Printf("NodeInvokeResult: %v (may not be implemented in mock)\n", err)
	} else {
		fmt.Println("Invoke result sent")
	}

	fmt.Println("\n=== Done ===")
}
