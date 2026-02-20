// Command discovery demonstrates mDNS gateway discovery on the local network.
//
// It uses DNS-SD to browse for _openclaw-gw._tcp services, which are
// advertised by running OpenClaw Gateway instances. Discovered gateways
// are printed with their connection details.
//
// Usage:
//
//	go run ./examples/discovery
//
// No mock server needed — this discovers real gateways on the LAN.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/discovery"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Gateway Discovery ===")
	fmt.Println()
	fmt.Println("Scanning for OpenClaw gateways on the local network...")
	fmt.Println("(Will wait up to 10 seconds)")
	fmt.Println()

	browser := discovery.NewBrowser()
	beacons, err := browser.Browse(ctx)
	if err != nil {
		log.Fatalf("Browse: %v", err)
	}

	if len(beacons) == 0 {
		fmt.Println("No gateways found on the local network.")
		fmt.Println()
		fmt.Println("Make sure an OpenClaw Gateway is running with mDNS enabled.")
		return
	}

	fmt.Printf("Found %d gateway(s):\n\n", len(beacons))

	for i, b := range beacons {
		fmt.Printf("  Gateway %d:\n", i+1)
		fmt.Printf("    Host:      %s\n", b.Host)
		fmt.Printf("    Port:      %d\n", b.Port)
		fmt.Printf("    Role:      %s\n", b.Role)
		fmt.Printf("    LanHost:   %s\n", b.LanHost)
		if b.GatewayTLS {
			fmt.Printf("    TLS:       yes (sha256: %s)\n", b.GatewayTLSFingerprint)
		} else {
			fmt.Printf("    TLS:       no\n")
		}
		if b.CanvasPort > 0 {
			fmt.Printf("    Canvas:    port %d\n", b.CanvasPort)
		}
		if b.CLIPath != "" {
			fmt.Printf("    CLI:       %s\n", b.CLIPath)
		}
		if b.TailnetDNS != "" {
			fmt.Printf("    Tailnet:   %s\n", b.TailnetDNS)
		}

		// Derive connection URLs using receiver methods.
		fmt.Printf("    WS URL:    %s\n", b.WebSocketURL())
		fmt.Printf("    HTTP URL:  %s\n", b.HTTPURL())
		fmt.Println()
	}

	fmt.Println("=== Done ===")
}
