package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// NodePairRequest requests pairing with the gateway.
func (c *Client) NodePairRequest(ctx context.Context, params protocol.NodePairRequestParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.pair.request", params)
}

// NodePairList lists pending node pairing requests.
func (c *Client) NodePairList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.pair.list", struct{}{})
}

// NodePairApprove approves a node pairing request.
func (c *Client) NodePairApprove(ctx context.Context, params protocol.NodePairApproveParams) error {
	return c.sendRPCVoid(ctx, "node.pair.approve", params)
}

// NodePairReject rejects a node pairing request.
func (c *Client) NodePairReject(ctx context.Context, params protocol.NodePairRejectParams) error {
	return c.sendRPCVoid(ctx, "node.pair.reject", params)
}

// NodePairVerify verifies a node pairing.
func (c *Client) NodePairVerify(ctx context.Context, params protocol.NodePairVerifyParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.pair.verify", params)
}
