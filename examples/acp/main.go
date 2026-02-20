// Command acp demonstrates an ACP (Agent Client Protocol) agent server.
//
// ACP is a JSON-RPC 2.0 over NDJSON protocol used by code editors to
// communicate with AI agents. This example implements a simple agent that:
//   - Responds to initialize with its capabilities
//   - Creates and manages sessions (new, load, fork, resume, list)
//   - Handles prompts by echoing them back
//   - Sends session update notifications
//   - Supports mode, model, and config option changes
//
// In production, this would bridge to an OpenClaw Gateway.
//
// Usage (pipe input):
//
//	echo '{"jsonrpc":"2.0","id":"1","method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | go run ./examples/acp
//
// Or use with an IDE that supports ACP.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/a3tai/openclaw-go/acp"
)

func main() {
	fmt.Fprintln(os.Stderr, "=== OpenClaw ACP Agent Example ===")
	fmt.Fprintln(os.Stderr, "Reading JSON-RPC requests from stdin, writing responses to stdout")
	fmt.Fprintln(os.Stderr, "")

	handler := &exampleHandler{
		sessions: make(map[string]*session),
	}

	srv := acp.NewServer(handler, os.Stdin, os.Stdout)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := srv.Serve(ctx); err != nil {
		log.Fatalf("Serve: %v", err)
	}

	fmt.Fprintln(os.Stderr, "Agent exiting")
}

// session tracks an active session.
type session struct {
	id    string
	mode  string
	model string
}

// exampleHandler implements acp.Handler.
type exampleHandler struct {
	sessions map[string]*session
	nextID   int
}

func (h *exampleHandler) Initialize(_ context.Context, req acp.InitializeRequest) (*acp.InitializeResponse, error) {
	name, version := "unknown", "unknown"
	if req.ClientInfo != nil {
		name = req.ClientInfo.Name
		version = req.ClientInfo.Version
	}
	fmt.Fprintf(os.Stderr, "[init] Client: %s %s, protocol: %d\n",
		name, version, req.ProtocolVersion)

	title := "OpenClaw Example Agent"
	return &acp.InitializeResponse{
		ProtocolVersion: acp.ProtocolVersion,
		AgentInfo: &acp.Implementation{
			Name:    "openclaw-example",
			Title:   &title,
			Version: "1.0.0",
		},
		AgentCapabilities: &acp.AgentCapabilities{
			LoadSession: true,
			PromptCapabilities: &acp.PromptCapabilities{
				Image:           true,
				EmbeddedContext: true,
			},
			SessionCapabilities: &acp.SessionCapabilities{
				List:   &acp.SessionListCapabilities{},
				Fork:   &acp.SessionForkCapabilities{},
				Resume: &acp.SessionResumeCapabilities{},
			},
		},
		AuthMethods: []acp.AuthMethod{
			{ID: "api_key", Name: "API Key"},
		},
	}, nil
}

func (h *exampleHandler) Authenticate(_ context.Context, req acp.AuthenticateRequest) (*acp.AuthenticateResponse, error) {
	fmt.Fprintf(os.Stderr, "[authenticate] Method: %s\n", req.MethodID)
	return &acp.AuthenticateResponse{}, nil
}

func (h *exampleHandler) NewSession(_ context.Context, req acp.NewSessionRequest) (*acp.NewSessionResponse, error) {
	h.nextID++
	id := fmt.Sprintf("sess-%d", h.nextID)
	h.sessions[id] = &session{id: id, mode: "code", model: "gpt-4"}
	fmt.Fprintf(os.Stderr, "[session/new] Created %s (cwd: %s)\n", id, req.CWD)
	return &acp.NewSessionResponse{
		SessionID: id,
		Modes: &acp.SessionModeState{
			AvailableModes: []acp.SessionMode{
				{ID: "code", Name: "Code"},
				{ID: "ask", Name: "Ask"},
			},
			CurrentModeID: "code",
		},
	}, nil
}

func (h *exampleHandler) LoadSession(_ context.Context, req acp.LoadSessionRequest) (*acp.LoadSessionResponse, error) {
	if _, ok := h.sessions[req.SessionID]; !ok {
		return nil, fmt.Errorf("session not found: %s", req.SessionID)
	}
	fmt.Fprintf(os.Stderr, "[session/load] Loaded %s\n", req.SessionID)
	return &acp.LoadSessionResponse{}, nil
}

func (h *exampleHandler) ListSessions(_ context.Context, _ acp.ListSessionsRequest) (*acp.ListSessionsResponse, error) {
	var sessions []acp.SessionInfo
	for id := range h.sessions {
		title := "Session " + id
		sessions = append(sessions, acp.SessionInfo{
			SessionID: id,
			CWD:       "/tmp",
			Title:     &title,
		})
	}
	return &acp.ListSessionsResponse{Sessions: sessions}, nil
}

func (h *exampleHandler) ForkSession(_ context.Context, req acp.ForkSessionRequest) (*acp.ForkSessionResponse, error) {
	if _, ok := h.sessions[req.SessionID]; !ok {
		return nil, fmt.Errorf("session not found: %s", req.SessionID)
	}
	h.nextID++
	id := fmt.Sprintf("sess-%d", h.nextID)
	h.sessions[id] = &session{id: id, mode: "code", model: "gpt-4"}
	fmt.Fprintf(os.Stderr, "[session/fork] Forked %s -> %s\n", req.SessionID, id)
	return &acp.ForkSessionResponse{SessionID: id}, nil
}

func (h *exampleHandler) ResumeSession(_ context.Context, req acp.ResumeSessionRequest) (*acp.ResumeSessionResponse, error) {
	if _, ok := h.sessions[req.SessionID]; !ok {
		return nil, fmt.Errorf("session not found: %s", req.SessionID)
	}
	fmt.Fprintf(os.Stderr, "[session/resume] Resumed %s\n", req.SessionID)
	return &acp.ResumeSessionResponse{}, nil
}

func (h *exampleHandler) Prompt(_ context.Context, req acp.PromptRequest) (*acp.PromptResponse, error) {
	fmt.Fprintf(os.Stderr, "[session/prompt] Session: %s, parts: %d\n",
		req.SessionID, len(req.Prompt))

	// Echo back the prompt content.
	for _, part := range req.Prompt {
		switch part.Type {
		case "text":
			fmt.Fprintf(os.Stderr, "  [prompt] Text: %s\n", truncate(part.Text, 100))
		case "image":
			fmt.Fprintf(os.Stderr, "  [prompt] Image: %s (%d bytes)\n", part.MimeType, len(part.Data))
		case "audio":
			fmt.Fprintf(os.Stderr, "  [prompt] Audio: %s (%d bytes)\n", part.MimeType, len(part.Data))
		case "resource_link":
			uri := ""
			if part.URI != nil {
				uri = *part.URI
			}
			fmt.Fprintf(os.Stderr, "  [prompt] Resource link: %s\n", uri)
		case "resource":
			if part.Resource != nil {
				fmt.Fprintf(os.Stderr, "  [prompt] Resource: %s\n", part.Resource.URI)
			}
		}
	}

	return &acp.PromptResponse{StopReason: "end_turn"}, nil
}

func (h *exampleHandler) Cancel(_ context.Context, req acp.CancelNotification) {
	fmt.Fprintf(os.Stderr, "[session/cancel] Session: %s\n", req.SessionID)
}

func (h *exampleHandler) SetSessionMode(_ context.Context, req acp.SetSessionModeRequest) (*acp.SetSessionModeResponse, error) {
	if s, ok := h.sessions[req.SessionID]; ok {
		s.mode = req.ModeID
	}
	fmt.Fprintf(os.Stderr, "[session/set_mode] Session: %s, modeId: %s\n",
		req.SessionID, req.ModeID)
	return &acp.SetSessionModeResponse{}, nil
}

func (h *exampleHandler) SetSessionModel(_ context.Context, req acp.SetSessionModelRequest) (*acp.SetSessionModelResponse, error) {
	if s, ok := h.sessions[req.SessionID]; ok {
		s.model = req.ModelID
	}
	fmt.Fprintf(os.Stderr, "[session/set_model] Session: %s, modelId: %s\n",
		req.SessionID, req.ModelID)
	return &acp.SetSessionModelResponse{}, nil
}

func (h *exampleHandler) SetSessionConfigOption(_ context.Context, req acp.SetSessionConfigOptionRequest) (*acp.SetSessionConfigOptionResponse, error) {
	fmt.Fprintf(os.Stderr, "[session/set_config_option] Session: %s, configId: %s, value: %s\n",
		req.SessionID, req.ConfigID, req.Value)
	return &acp.SetSessionConfigOptionResponse{
		ConfigOptions: []acp.SessionConfigOption{
			{Type: "select", ID: req.ConfigID, Name: req.ConfigID, CurrentValue: req.Value},
		},
	}, nil
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
