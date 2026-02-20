package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// ConfigGet retrieves the current gateway configuration.
func (c *Client) ConfigGet(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "config.get", protocol.ConfigGetParams{})
}

// ConfigSet replaces the gateway configuration.
func (c *Client) ConfigSet(ctx context.Context, params protocol.ConfigSetParams) error {
	return c.sendRPCVoid(ctx, "config.set", params)
}

// ConfigApply applies a configuration change with optional restart.
func (c *Client) ConfigApply(ctx context.Context, params protocol.ConfigApplyParams) error {
	return c.sendRPCVoid(ctx, "config.apply", params)
}

// ConfigPatch patches the gateway configuration.
func (c *Client) ConfigPatch(ctx context.Context, params protocol.ConfigPatchParams) error {
	return c.sendRPCVoid(ctx, "config.patch", params)
}

// ConfigSchema retrieves the configuration JSON schema.
func (c *Client) ConfigSchema(ctx context.Context) (*protocol.ConfigSchemaResponse, error) {
	var result protocol.ConfigSchemaResponse
	if err := c.sendRPCTyped(ctx, "config.schema", struct{}{}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
