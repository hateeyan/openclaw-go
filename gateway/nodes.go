package gateway

import (
	"context"
	"encoding/json"

	"github.com/a3tai/openclaw-go/protocol"
)

// NodeList lists connected nodes.
func (c *Client) NodeList(ctx context.Context) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.list", struct{}{})
}

// NodeDescribe describes a specific node.
func (c *Client) NodeDescribe(ctx context.Context, params protocol.NodeDescribeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.describe", params)
}

// NodeInvoke invokes a command on a node.
func (c *Client) NodeInvoke(ctx context.Context, params protocol.NodeInvokeParams) (json.RawMessage, error) {
	return c.sendRPC(ctx, "node.invoke", params)
}

// NodeInvokeResult sends an invoke result from a node to the gateway.
func (c *Client) NodeInvokeResult(ctx context.Context, params protocol.NodeInvokeResultParams) error {
	return c.sendRPCVoid(ctx, "node.invoke.result", params)
}

// NodeEvent sends an event from a node to the gateway.
func (c *Client) NodeEvent(ctx context.Context, params protocol.NodeEventParams) error {
	return c.sendRPCVoid(ctx, "node.event", params)
}

// NodeRename renames a node.
func (c *Client) NodeRename(ctx context.Context, params protocol.NodeRenameParams) error {
	return c.sendRPCVoid(ctx, "node.rename", params)
}
