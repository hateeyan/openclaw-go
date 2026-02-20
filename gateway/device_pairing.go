package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// DevicePairList lists device pairing entries.
func (c *Client) DevicePairList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "device.pair.list", struct{}{})
}

// DevicePairApprove approves a device pairing request.
func (c *Client) DevicePairApprove(ctx context.Context, params protocol.DevicePairApproveParams) error {
	return c.sendRPCVoid(ctx, "device.pair.approve", params)
}

// DevicePairReject rejects a device pairing request.
func (c *Client) DevicePairReject(ctx context.Context, params protocol.DevicePairRejectParams) error {
	return c.sendRPCVoid(ctx, "device.pair.reject", params)
}

// DevicePairRemove removes a paired device.
func (c *Client) DevicePairRemove(ctx context.Context, params protocol.DevicePairRemoveParams) error {
	return c.sendRPCVoid(ctx, "device.pair.remove", params)
}

// DeviceTokenRotate rotates a device token.
func (c *Client) DeviceTokenRotate(ctx context.Context, params protocol.DeviceTokenRotateParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "device.token.rotate", params)
}

// DeviceTokenRevoke revokes a device token.
func (c *Client) DeviceTokenRevoke(ctx context.Context, params protocol.DeviceTokenRevokeParams) error {
	return c.sendRPCVoid(ctx, "device.token.revoke", params)
}
