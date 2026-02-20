# gateway

```
import "github.com/a3tai/openclaw-go/gateway"
```

Package `gateway` implements a WebSocket client for the OpenClaw Gateway protocol. It handles the full connection lifecycle including the challenge/connect/hello-ok handshake, background read loop with frame dispatch, keepalive pings, and pending request/response correlation.

## Creating a Client

```go
client := gateway.NewClient(
    gateway.WithToken("my-token"),
    gateway.WithRole(protocol.RoleOperator),
    gateway.WithScopes(protocol.ScopeOperatorRead, protocol.ScopeOperatorWrite),
    gateway.WithOnEvent(func(ev protocol.Event) {
        fmt.Printf("event: %s\n", ev.EventName)
    }),
)
defer client.Close()

ctx := context.Background()
if err := client.Connect(ctx, "ws://localhost:18789/ws"); err != nil {
    log.Fatal(err)
}

hello := client.Hello()
fmt.Printf("protocol=%d server=%s\n", hello.Protocol, hello.Server.Version)
```

## Configuration Options

| Option | Description |
|--------|-------------|
| `WithToken(token)` | Bearer token authentication |
| `WithPassword(password)` | Password authentication |
| `WithClientInfo(info)` | Override default client identity |
| `WithRole(role)` | Connection role: `RoleOperator` or `RoleNode` |
| `WithScopes(scopes...)` | Operator capability scopes |
| `WithCaps(caps...)` | Node capability categories |
| `WithCommands(commands...)` | Node command allowlist |
| `WithPermissions(perms)` | Node permission toggles |
| `WithDevice(device)` | Device identity for pairing |
| `WithLocale(locale)` | Locale string (default: `"en-US"`) |
| `WithUserAgent(ua)` | User-agent string |
| `WithTLSConfig(cfg)` | Custom TLS configuration |
| `WithConnectTimeout(d)` | Handshake timeout (default: 10s) |
| `WithOnEvent(fn)` | Callback for incoming events |
| `WithOnInvoke(fn)` | Handler for incoming invocations (node mode) |

## Core Methods

```go
// Connection lifecycle
err := client.Connect(ctx, wsURL)
err := client.Close()
<-client.Done()  // closed when client shuts down
hello := client.Hello()  // server hello-ok payload

// Low-level
resp, err := client.Send(ctx, method, params)
err := client.SendEvent(eventName, payload)
```

## RPC Methods

The client provides typed convenience methods for all 96+ gateway RPCs. Each method marshals the params, sends the request, waits for the response, and unmarshals the typed result.

### Chat

```go
result, err := client.ChatSend(ctx, protocol.ChatSendParams{
    SessionKey: "main",
    Message:    "Hello!",
})
history, err := client.ChatHistory(ctx, protocol.ChatHistoryParams{SessionKey: "main"})
err := client.ChatAbort(ctx, protocol.ChatAbortParams{SessionKey: "main"})
err := client.ChatInject(ctx, protocol.ChatInjectParams{SessionKey: "main", Messages: msgs})
```

### Agent

```go
agent, err := client.Agent(ctx, protocol.AgentParams{})
identity, err := client.AgentIdentity(ctx, protocol.AgentIdentityParams{})
result, err := client.AgentWait(ctx, protocol.AgentWaitParams{})
```

### Sessions

```go
sessions, err := client.SessionsList(ctx, protocol.SessionsListParams{})
preview, err := client.SessionsPreview(ctx, protocol.SessionsPreviewParams{SessionKey: key})
err := client.SessionsPatch(ctx, protocol.SessionsPatchParams{SessionKey: key, Label: &label})
usage, err := client.SessionsUsage(ctx, protocol.SessionsUsageParams{SessionKey: key})
err := client.SessionsReset(ctx, protocol.SessionsResetParams{SessionKey: key})
err := client.SessionsDelete(ctx, protocol.SessionsDeleteParams{SessionKey: key})
```

### Agents CRUD

```go
list, err := client.AgentsList(ctx)
created, err := client.AgentsCreate(ctx, protocol.AgentsCreateParams{Name: "my-agent"})
err := client.AgentsUpdate(ctx, protocol.AgentsUpdateParams{AgentID: id})
err := client.AgentsDelete(ctx, protocol.AgentsDeleteParams{AgentID: id})

// Agent files
files, err := client.AgentsFilesList(ctx, protocol.AgentsFilesListParams{AgentID: id})
content, err := client.AgentsFilesGet(ctx, protocol.AgentsFilesGetParams{AgentID: id, Path: path})
err := client.AgentsFilesSet(ctx, protocol.AgentsFilesSetParams{AgentID: id, Path: path, Content: data})
```

### Config

```go
cfg, err := client.ConfigGet(ctx, protocol.ConfigGetParams{})
schema, err := client.ConfigSchema(ctx)
err := client.ConfigSet(ctx, protocol.ConfigSetParams{Key: "key", Value: val})
err := client.ConfigPatch(ctx, protocol.ConfigPatchParams{Patch: patch})
err := client.ConfigApply(ctx, protocol.ConfigApplyParams{Restart: true})
```

### Exec Approvals

```go
result, err := client.ExecApprovalRequest(ctx, protocol.ExecApprovalRequestParams{...})
result, err := client.ResolveExecApproval(ctx, protocol.ExecApprovalResolveParams{
    ID: "approval-id", Decision: "approved",
})
result, err := client.ExecApprovalWaitDecision(ctx, protocol.ExecApprovalWaitDecisionParams{ID: id})
```

### Nodes & Pairing

```go
// Node operations
result, err := client.NodeList(ctx)
desc, err := client.NodeDescribe(ctx, protocol.NodeDescribeParams{NodeID: id})
result, err := client.NodeInvoke(ctx, protocol.NodeInvokeParams{NodeID: id, Cap: "tool", Args: args})

// Node pairing
result, err := client.NodePairRequest(ctx, protocol.NodePairRequestParams{...})
pairs, err := client.NodePairList(ctx)
err := client.NodePairApprove(ctx, protocol.NodePairApproveParams{ID: id})
err := client.NodePairReject(ctx, protocol.NodePairRejectParams{ID: id})

// Device pairing
pairs, err := client.DevicePairList(ctx)
err := client.DevicePairApprove(ctx, protocol.DevicePairApproveParams{ID: id})
err := client.DevicePairRemove(ctx, protocol.DevicePairRemoveParams{DeviceID: id})
result, err := client.DeviceTokenRotate(ctx, protocol.DeviceTokenRotateParams{DeviceID: id})
```

### Cron

```go
jobs, err := client.CronList(ctx, protocol.CronListParams{})
status, err := client.CronStatus(ctx)
err := client.CronAdd(ctx, protocol.CronAddParams{Name: "nightly", Schedule: sched})
err := client.CronUpdate(ctx, protocol.CronUpdateParams{ID: id})
err := client.CronRemove(ctx, protocol.CronRemoveParams{ID: id})
err := client.CronRun(ctx, protocol.CronRunParams{ID: id})
runs, err := client.CronRuns(ctx, protocol.CronRunsParams{ID: id})
```

### TTS

```go
status, err := client.TTSStatus(ctx)
providers, err := client.TTSProviders(ctx)
err := client.TTSEnable(ctx)
err := client.TTSDisable(ctx)
result, err := client.TTSConvert(ctx, protocol.TTSConvertParams{Text: "Hello"})
err := client.TTSSetProvider(ctx, protocol.TTSSetProviderParams{Provider: "eleven"})
```

### Other Methods

```go
// Health & status
health, err := client.Health(ctx)
status, err := client.Status(ctx)
presence, err := client.Presence(ctx)

// Models
models, err := client.ModelsList(ctx)

// Logs
logs, err := client.LogsTail(ctx, protocol.LogsTailParams{Lines: 100})

// Channels & talk
chStatus, err := client.ChannelsStatus(ctx, protocol.ChannelsStatusParams{})
talkCfg, err := client.TalkConfig(ctx, protocol.TalkConfigParams{})

// Skills
skills, err := client.SkillsStatus(ctx, protocol.SkillsStatusParams{})
bins, err := client.SkillsBins(ctx)

// Wizard
wizard, err := client.WizardStart(ctx, protocol.WizardStartParams{})
next, err := client.WizardNext(ctx, protocol.WizardNextParams{})

// Misc
err := client.Wake(ctx, protocol.WakeParams{})
err := client.SendMessage(ctx, protocol.SendParams{...})
```

## Node Mode

To connect as a capability node, use `WithRole(protocol.RoleNode)` and `WithOnInvoke`:

```go
client := gateway.NewClient(
    gateway.WithToken("node-token"),
    gateway.WithRole(protocol.RoleNode),
    gateway.WithCaps("search", "calculator"),
    gateway.WithOnInvoke(func(inv protocol.Invoke) protocol.InvokeResponse {
        // Handle the invocation
        result := processInvocation(inv.Cap, inv.Args)
        return protocol.InvokeResponse{
            OK:      true,
            Payload: result,
        }
    }),
)
```

## Concurrency

The client is safe for concurrent use. Writes are mutex-protected, request IDs use an atomic counter, and the pending response map is synchronized. The background read loop and tick loop run in separate goroutines.
