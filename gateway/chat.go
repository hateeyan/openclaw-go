package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ChatSend sends a chat message and returns the chat event response.
func (c *Client) ChatSend(ctx context.Context, params protocol.ChatSendParams) (*protocol.ChatEvent, error) {
	var ev protocol.ChatEvent
	if err := c.sendRPCTyped(ctx, "chat.send", params, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

// ChatHistory retrieves chat history for a session.
func (c *Client) ChatHistory(ctx context.Context, params protocol.ChatHistoryParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "chat.history", params)
}

// ChatAbort aborts a running chat session.
func (c *Client) ChatAbort(ctx context.Context, params protocol.ChatAbortParams) error {
	return c.sendRPCVoid(ctx, "chat.abort", params)
}

// ChatInject injects a message into a chat session.
func (c *Client) ChatInject(ctx context.Context, params protocol.ChatInjectParams) error {
	return c.sendRPCVoid(ctx, "chat.inject", params)
}
