package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// UpdateRun triggers a gateway update run.
func (c *Client) UpdateRun(ctx context.Context, params protocol.UpdateRunParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "update.run", params)
}

// PushTest sends a test push notification.
func (c *Client) PushTest(ctx context.Context, params protocol.PushTestParams) (*protocol.PushTestResult, error) {
	var result protocol.PushTestResult
	if err := c.sendRPCTyped(ctx, "push.test", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BrowserRequest makes a browser request via the gateway.
func (c *Client) BrowserRequest(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, "browser.request", params)
}

// VoiceWakeGet retrieves the voice wake configuration.
func (c *Client) VoiceWakeGet(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "voicewake.get", struct{}{})
}

// VoiceWakeSet sets the voice wake configuration.
func (c *Client) VoiceWakeSet(ctx context.Context, params any) error {
	return c.sendRPCVoid(ctx, "voicewake.set", params)
}

// UsageStatus retrieves usage status.
func (c *Client) UsageStatus(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "usage.status", struct{}{})
}

// UsageCost retrieves usage cost information.
func (c *Client) UsageCost(ctx context.Context, params any) (json.RawMessage, error) {
	return c.sendRPC(ctx, "usage.cost", params)
}

// Poll creates a poll.
func (c *Client) Poll(ctx context.Context, params protocol.PollParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "poll", params)
}
