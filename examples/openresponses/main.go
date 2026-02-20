// Command openresponses demonstrates the OpenAI Responses API via OpenClaw.
//
// This is the newer API (vs Chat Completions) that supports:
//   - Structured input items (messages, function call outputs, images)
//   - SSE streaming with typed events
//   - Tool/function calling
//
// Usage:
//
//	go run ./examples/openresponses
//
// Requires a gateway with the /v1/responses endpoint.
// Note: The mock server does not implement this endpoint yet.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/a3tai/openclaw-go/openresponses"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== OpenClaw Open Responses Example ===")
	fmt.Println()

	client := &openresponses.Client{
		BaseURL:    "http://localhost:18789",
		Token:      "example-token",
		AgentID:    "main",
		SessionKey: "example-session",
	}

	// --- Non-streaming request ---
	fmt.Println("--- Non-streaming ---")
	demonstrateCreate(ctx, client)

	// --- Streaming request ---
	fmt.Println("\n--- Streaming ---")
	demonstrateStream(ctx, client)

	fmt.Println("\n=== Done ===")
}

func demonstrateCreate(ctx context.Context, client *openresponses.Client) {
	maxTokens := 200
	resp, err := client.Create(ctx, openresponses.Request{
		Model: "openclaw:main",
		Input: openresponses.InputFromItems([]openresponses.InputItem{
			openresponses.MessageItem("user", "Explain what OpenClaw is in one sentence."),
		}),
		MaxOutputTokens: &maxTokens,
	})
	if err != nil {
		log.Printf("Create: %v (endpoint may not be available)", err)
		return
	}

	fmt.Printf("Response ID: %s\n", resp.ID)
	fmt.Printf("Status: %s\n", resp.Status)
	for _, item := range resp.Output {
		if item.Type == "message" {
			for _, part := range item.Content {
				if part.Type == "output_text" {
					fmt.Printf("Text: %s\n", part.Text)
				}
			}
		}
	}
	fmt.Printf("Usage: input=%d, output=%d, total=%d\n",
		resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.TotalTokens)
}

func demonstrateStream(ctx context.Context, client *openresponses.Client) {
	// Build a request with a function tool.
	req := openresponses.Request{
		Model: "openclaw:main",
		Input: openresponses.InputFromItems([]openresponses.InputItem{
			openresponses.MessageItem("user", "What is the weather in San Francisco?"),
		}),
		Tools: []openresponses.ToolDefinition{
			{
				Type: "function",
				Function: openresponses.FunctionTool{
					Name:        "get_weather",
					Description: "Get the current weather for a location",
					Parameters: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"location": map[string]string{
								"type":        "string",
								"description": "City name",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
	}

	stream, err := client.CreateStream(ctx, req)
	if err != nil {
		log.Printf("CreateStream: %v (endpoint may not be available)", err)
		return
	}
	defer stream.Close()

	fmt.Print("Streaming: ")
	for {
		ev, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Recv: %v", err)
			break
		}

		switch ev.EventType {
		case "response.created":
			fmt.Print("[created] ")
		case "response.output_text.delta":
			fmt.Print(".")
		case "response.output_text.done":
			fmt.Print(" [text done] ")
		case "response.function_call_arguments.delta":
			fmt.Print("f")
		case "response.completed":
			fmt.Print("[completed]")
		}
	}
	fmt.Println()
}

// Ensure time import is used (for context timeout).
var _ = time.Second
