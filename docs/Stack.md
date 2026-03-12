# Stack

## Purpose
Build a cross-platform Go CLI that implements the Trello command contract defined in the command specification.

This is not just a thin wrapper around Trello anymore. It is a stable, machine-oriented product surface for AI agents and terminal users.

That means the stack must optimize for:

- deterministic JSON output
- stable validation and error semantics
- safe credential handling
- clear command composition
- low dependency count
- single-binary cross-platform distribution

---

## Design Constraints from the Command Spec

The command spec changes the architectural shape in a few important ways:

- all commands return JSON on stdout
- success and error envelopes are standardized
- interactive browser login is a first-class auth flow
- manual credential entry is also supported
- attachments require local file handling and multipart upload support
- command groups are broad enough that a real command tree is required
- output consistency matters more than convenience shortcuts

So the stack should be built around contract enforcement, not around a loose collection of HTTP calls.

---

## Core Technical Decisions

### 1. Language
**Go**

Why:

- compiles to a single self-contained binary
- strong Windows, macOS, and Linux support
- fast startup time for short-lived CLI invocations
- standard library is enough for HTTP, JSON, timeouts, multipart uploads, context, and logging
- easy cross-compilation and CI packaging

Decision:
**Use Go.**

---

## CLI Framework

### 2. Cobra
**Role:** command structure, flags, help text, subcommands, shared validation hooks

Use Cobra as the command framework because the command surface is now clearly a resource-oriented command tree, not a flat utility.

Representative shape:

```text
trello
  auth
    status
    login
    set
    clear
  boards
    list
    get
  lists
    list
    create
    update
    archive
    move
  cards
    list
    get
    create
    update
    move
    archive
    delete
  comments
    list
    add
    update
    delete
  checklists
    list
    create
    delete
    items
      add
      update
      delete
  attachments
    list
    add-file
    add-url
    delete
  labels
    list
    create
    add
    remove
  members
    list
    add
    remove
  search
    cards
    boards
  version
```

Why Cobra fits:

- nested command groups are a first-class need
- shared persistent flags belong at the root or group level
- per-command validation hooks are easy to keep consistent
- help text and shell completion come for free
- command wiring stays organized instead of turning into flag spaghetti

Why not `flag` alone:

- fine for tiny one-command tools
- bad once nested subcommands and shared behavior appear
- would push too much contract logic into hand-rolled plumbing

Why not `urfave/cli` as the default:

- good for flatter CLIs
- less compelling once the command tree gets deep and resource-oriented
- Cobra is the more natural fit for this spec

Decision:
**Use Cobra as the command framework.**

---

### 3. Viper
**Role:** non-secret configuration, environment binding, defaults, optional profiles

Use Viper, but give it a strict job description.

Viper should manage:

- non-secret config values
- environment variable binding
- config file loading and precedence
- defaults for timeout, retry policy, pretty-printing, and profile selection

Viper should **not** be the default secret store.

Do not use Viper as a global mutable junk drawer. That is how you get configuration that feels haunted.

Recommended precedence:

1. explicit flags
2. environment variables
3. config file
4. defaults

Recommended Viper-managed settings:

- `timeout`
- `max_retries`
- `retry_mutations`
- `pretty`
- `profile`
- optional default board or member identifiers only if they prove useful later

Decision:
**Use Viper, but only for non-secret config and env binding.**

---

## Credentials and Authentication

### 4. Credential storage abstraction
**Role:** safe persistence of Trello API key and token outside normal config files

The command spec requires both:

- `trello auth login`
- `trello auth set --api-key ... --token ...`

That means credentials must persist across future commands.

Design rule:
Create a credential storage abstraction rather than coupling secrets to Viper config.

Responsibilities:

- store API key and token for the active profile
- load credentials for command execution
- clear credentials on `trello auth clear`
- avoid printing or logging secrets
- support cross-platform implementations
- support headless fallback behavior for CI or automation when interactive login is not practical

Implementation note:
Prefer an OS-backed secure store when available. If a fallback is needed for specific environments, keep it explicit and opt-in rather than silently dropping secrets into plain config files.

Decision:
**Add a dedicated credential store layer. Do not treat Viper config as the primary secret store.**

---

### 5. Auth subsystem
**Role:** interactive login, manual credential setup, status verification, and auth-mode reporting

Auth is now its own subsystem, not just “append key and token to every request.”

Required V1 flows:

- `trello auth status`
- `trello auth login`
- `trello auth set`
- `trello auth clear`

Auth subsystem responsibilities:

- launch or guide the user through Trello authorization
- capture or accept credentials
- verify credentials via `members/me`
- persist validated credentials
- report `authMode` consistently as `interactive`, `manual`, or `null`
- fail protected commands with `AUTH_REQUIRED` when no valid auth exists
- map Trello auth failures to `AUTH_INVALID`

Design rule:
Keep auth verification centralized so each command does not invent its own version of “am I logged in?”

Decision:
**Implement auth as a dedicated internal subsystem.**

---

## Networking Layer

### 6. `net/http` + `encoding/json`
**Role:** Trello API transport and serialization

Use the Go standard library for HTTP and JSON.

This tool does not need a heavy third-party REST client. `net/http` plus `encoding/json` is enough for:

- GET/POST/PUT/DELETE requests
- query parameters and headers
- timeouts and context cancellation
- JSON request and response handling
- multipart file uploads for attachments
- retry handling

Design rule:
Create a small internal Trello client package instead of scattering raw HTTP calls across Cobra commands.

Trello client responsibilities:

- build request URLs and query params
- inject auth credentials safely
- execute requests with context and timeout control
- decode success responses
- normalize Trello and transport errors
- support multipart file uploads for `attachments add-file`
- map rate-limit responses into retry logic and CLI errors

Decision:
**Use `net/http` and `encoding/json`, not a third-party REST abstraction.**

---

## Output Contract

### 7. JSON-only stdout
**Role:** stable machine interface for agents and automation

The command spec is explicit: all commands return JSON on stdout.

Success envelope:

```json
{
  "ok": true,
  "data": {}
}
```

Error envelope:

```json
{
  "ok": false,
  "error": {
    "code": "STRING_CODE",
    "message": "Human-readable message"
  }
}
```

Rules:

- stdout is reserved for JSON only
- stderr is reserved for debug or incidental diagnostics only
- do not emit decorative tables, plain-text summaries, color, spinners, or progress bars in normal command execution
- `--pretty` may control pretty-printed JSON only
- do not ship a text output mode in V1

Decision:
**Make JSON the only supported stdout contract in V1.**

---

## Contract and Validation Layer

### 8. Internal contract package
**Role:** enforce response envelopes, standard error codes, and reusable validation behavior

The command spec defines stable external behavior. Do not let each command invent its own JSON and error mapping.

Create a small internal contract layer that owns:

- success envelope creation
- error envelope creation
- standard error code constants
- flag and argument validation helpers
- Trello-to-CLI error normalization
- shared response-field shaping where needed

Standard error codes from the spec:

- `AUTH_REQUIRED`
- `AUTH_INVALID`
- `NOT_FOUND`
- `VALIDATION_ERROR`
- `CONFLICT`
- `RATE_LIMITED`
- `HTTP_ERROR`
- `FILE_NOT_FOUND`
- `UNSUPPORTED`
- `UNKNOWN_ERROR`

Design rule:
Cobra commands should call reusable validators and return normalized contract errors, not hand-roll string messages ad hoc.

Decision:
**Add an internal contract/validation layer.**

---

## Error Handling

### 9. Predictable failures
CLI failures should be boring, stable, and machine-readable.

Requirements:

- return the standard JSON error envelope on stdout
- use non-zero exit codes on failure
- keep human-readable messages concise and safe
- normalize Trello HTTP failures into spec-defined codes
- never leak secrets in logs or errors

Suggested mapping examples:

- missing required flags -> `VALIDATION_ERROR`
- missing credentials -> `AUTH_REQUIRED`
- invalid Trello key/token -> `AUTH_INVALID`
- missing Trello resource -> `NOT_FOUND`
- local file path does not exist -> `FILE_NOT_FOUND`
- unsupported/excluded feature -> `UNSUPPORTED`
- HTTP 429 -> `RATE_LIMITED`
- unexpected transport or upstream failures -> `HTTP_ERROR` or `UNKNOWN_ERROR`

Decision:
**Keep failure semantics centralized and deterministic.**

---

## Rate Limiting and Retries

### 10. Backoff policy
Trello rate limits are real, and pretending otherwise is a cute way to build a flaky tool.

Requirements:

- detect HTTP 429 responses
- use bounded exponential backoff with jitter
- respect `context.Context` cancellation
- cap retry attempts
- keep retry logic in the API client layer, not in commands

Recommended initial policy:

- retry idempotent reads by default
- do not automatically retry unsafe mutations unless explicitly enabled
- make retry behavior configurable through Viper-managed settings and flags

Suggested settings:

- `timeout`
- `max_retries`
- `retry_mutations`

Decision:
**Implement retry policy in the Trello client layer, not in command handlers.**

---

## Logging and Observability

### 11. Keep logging simple
For V1, normal command output belongs on stdout as JSON and nowhere else.

Logging rules:

- send diagnostics to stderr only
- gate request/response diagnostics behind `--verbose`
- never log secrets or full credential-bearing URLs
- prefer stdlib logging unless a real structured-logging need appears later

Do not add a logging framework just because enterprise software enjoys wearing too many belts.

Decision:
**Use simple logging with strict secret redaction.**

---

## Command-Surface Guidance

### 12. Resource-oriented command design
The command spec already defines the resource surface, so the stack should preserve that shape instead of introducing escape hatches that undermine it.

Guidelines:

- nouns at the first level (`boards`, `cards`, `lists`, `labels`, `members`)
- verbs or actions at the second level (`list`, `get`, `create`, `update`, `delete`)
- nested subcommands only where the domain clearly needs them (`checklists items ...`)
- flag names must match the command spec exactly where the spec defines them

Important change from the earlier stack:

- do **not** treat `raw request` as part of the stable V1 surface

A raw passthrough command would punch a hole through the contract by allowing inconsistent behavior, unsupported fields, and undocumented semantics. If a raw command ever appears later, it should be hidden or explicitly marked experimental.

Decision:
**Build the command tree around the spec; do not ship a stable `raw` escape hatch in V1.**

---

## Configuration Layout

### 13. Minimal non-secret config
If config file support is enabled, keep it small and boring.

Example:

```yaml
profile: default
pretty: false
timeout: 15s
max_retries: 3
retry_mutations: false

profiles:
  default:
    auth_mode: interactive
```

Rules:

- config stores non-secret preferences, not default plaintext credentials
- environment variables remain useful for CI and headless use
- profile support should be optional and narrow, not a tiny configuration religion

Decision:
**Keep config minimal and keep secrets out of it by default.**

---

## Suggested Internal Package Layout

```text
/cmd
  /trello
    main.go
    root.go
    auth.go
    boards.go
    lists.go
    cards.go
    comments.go
    checklists.go
    attachments.go
    labels.go
    members.go
    search.go
    version.go

/internal
  /contract
    response.go
    errors.go
    validation.go
  /config
    config.go
  /credentials
    store.go
  /auth
    login.go
    status.go
    set.go
    clear.go
  /trello
    client.go
    auth.go
    boards.go
    lists.go
    cards.go
    comments.go
    checklists.go
    attachments.go
    labels.go
    members.go
    search.go
    errors.go
    retry.go
  /version
    version.go
```

Notes:

- `/cmd/trello` owns Cobra wiring and minimal orchestration
- `/internal/contract` owns the external command contract behavior
- `/internal/config` owns Viper setup and precedence
- `/internal/credentials` owns secret persistence abstraction
- `/internal/auth` owns login and auth-state workflows
- `/internal/trello` owns API concerns and resource-specific client logic

This keeps command handlers thin and keeps the HTTP and contract logic reusable.

---

## Dependencies

### Required

- `github.com/spf13/cobra`
- `github.com/spf13/viper`

### Standard library

- `net/http`
- `mime/multipart`
- `encoding/json`
- `context`
- `time`
- `os`
- `fmt`
- `errors`
- `path/filepath`
- `log/slog` or `log`

### Additional implementation dependency
Select a small credential-store library only if needed to back the credential abstraction on each platform.

Rule:
Do not import five helper libraries for a tool whose real job is still “call Trello, print JSON, leave.”

---

## Build and Release

### 14. Distribution model
Ship as a single compiled binary for:

- Linux
- macOS
- Windows

Release requirements:

- embed version, commit, and build date
- support `trello version`
- publish checksums
- test cross-platform builds in CI

Suggested build metadata variables:

- version
- commit
- date

Decision:
**Single-binary distribution remains the right deployment model.**

---

## Testing Strategy

### 15. Test layers

#### Unit tests
Focus on:

- request construction
- query and path parameter handling
- JSON decoding and envelope creation
- config precedence
- auth-state validation
- credential-store integration boundaries
- error normalization
- retry behavior
- multipart attachment request construction

#### Integration tests
Use:

- mock HTTP servers
- optional real Trello sandbox tests behind environment-gated execution

Golden rule:
Do not make the whole test suite depend on live Trello. That is not a test plan; that is weather with extra steps.

---

## Security Rules

### 16. Secret handling

- never print API keys or tokens
- redact secrets in logs and error messages
- never include credential-bearing URLs in diagnostics
- avoid plain-text disk storage by default
- verify credentials before claiming auth is configured
- document scope expectations for interactive and manual auth modes

Decision:
**Secret handling is a first-order design concern, not an afterthought.**

---

## Recommended V1 Stack Summary

### Final stack

- **Language:** Go
- **CLI framework:** Cobra
- **Config:** Viper for non-secret settings and env binding
- **Credential storage:** dedicated credential-store abstraction
- **Auth subsystem:** interactive login + manual credential setup + status + clear
- **HTTP client:** `net/http`
- **Serialization:** `encoding/json`
- **Uploads:** stdlib multipart support
- **Contract layer:** standardized JSON envelopes and validation helpers
- **Logging:** stdlib logging with redaction
- **Build output:** single binary
- **Default output mode:** JSON-only on stdout

### Why this stack
Because the tool is now a stable agent-facing interface, not just a convenience wrapper.

That means the stack should optimize for:

- reliability
- portability
- strict contract behavior
- safe authentication
- low startup overhead
- low dependency count
- machine-readable output

Not for framework tourism.

---

## Non-goals for V1

Do not add these unless a real need appears:

- plugin architecture
- local database
- interactive TUI beyond the auth flow
- generic raw endpoint passthrough as a stable feature
- websocket/event streaming abstractions
- giant SDK layer wrapping every Trello endpoint immediately
- stateful background daemon

That is how a clean CLI turns into a haunted Victorian mansion with an API key in the attic.
