package gateway

import (
	"context"

	"github.com/a3tai/openclaw-go/protocol"
)

// ResolveExecApproval resolves a pending exec approval request.
// Requires the operator.approvals scope.
func (c *Client) ResolveExecApproval(ctx context.Context, params protocol.ExecApprovalResolveParams) (*protocol.ExecApprovalResolveResult, error) {
	var result protocol.ExecApprovalResolveResult
	if err := c.sendRPCTyped(ctx, "exec.approval.resolve", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalRequest submits a new exec approval request.
func (c *Client) ExecApprovalRequest(ctx context.Context, params protocol.ExecApprovalRequestParams) (*protocol.ExecApprovalRequestResult, error) {
	var result protocol.ExecApprovalRequestResult
	if err := c.sendRPCTyped(ctx, "exec.approval.request", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalWaitDecision waits for a decision on a pending exec approval.
// Requires the operator.approvals scope.
func (c *Client) ExecApprovalWaitDecision(ctx context.Context, params protocol.ExecApprovalWaitDecisionParams) (*protocol.ExecApprovalWaitDecisionResult, error) {
	var result protocol.ExecApprovalWaitDecisionResult
	if err := c.sendRPCTyped(ctx, "exec.approval.waitDecision", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecApprovalsGet retrieves the exec approvals configuration.
func (c *Client) ExecApprovalsGet(ctx context.Context) (*protocol.ExecApprovalsSnapshot, error) {
	var snap protocol.ExecApprovalsSnapshot
	if err := c.sendRPCTyped(ctx, "exec.approvals.get", protocol.ExecApprovalsGetParams{}, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// ExecApprovalsSet updates the exec approvals configuration.
func (c *Client) ExecApprovalsSet(ctx context.Context, params protocol.ExecApprovalsSetParams) error {
	return c.sendRPCVoid(ctx, "exec.approvals.set", params)
}

// ExecApprovalsNodeGet retrieves exec approvals for a specific node.
func (c *Client) ExecApprovalsNodeGet(ctx context.Context, params protocol.ExecApprovalsNodeGetParams) (*protocol.ExecApprovalsSnapshot, error) {
	var snap protocol.ExecApprovalsSnapshot
	if err := c.sendRPCTyped(ctx, "exec.approvals.node.get", params, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// ExecApprovalsNodeSet updates exec approvals for a specific node.
func (c *Client) ExecApprovalsNodeSet(ctx context.Context, params protocol.ExecApprovalsNodeSetParams) error {
	return c.sendRPCVoid(ctx, "exec.approvals.node.set", params)
}
