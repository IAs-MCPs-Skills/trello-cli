# Trello CLI

[![CI](https://github.com/Scale-Flow/trello-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/Scale-Flow/trello-cli/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod-go-version/Scale-Flow/trello-cli)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/Scale-Flow/trello-cli)](https://github.com/Scale-Flow/trello-cli/releases)

A cross-platform Go CLI for Trello with machine-friendly JSON output.

Every command writes a JSON envelope to `stdout`, which makes the CLI useful both for terminal users and for automation that needs predictable responses.

## Highlights

- Stable JSON success and error envelopes
- Interactive browser login and manual credential setup
- Broad Trello command coverage for boards, lists, cards, comments, labels, members, checklists, attachments, custom fields, and search
- Single-binary Go CLI with minimal runtime requirements

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap Scale-Flow/tap
brew install trello-cli
```

### Go Install

```bash
go install github.com/Scale-Flow/trello-cli/cmd/trello@latest
```

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/Scale-Flow/trello-cli/releases).

## Prerequisites: Trello API Key

Before using the CLI, you need a Trello API key and token:

1. Log in to [Trello](https://trello.com)
2. Go to the [Power-Up admin page](https://trello.com/power-ups/admin)
3. Create a new Power-Up (it doesn't need to do anything — it's just the container for your API credentials)
4. In your Power-Up's settings, go to the **API Key** tab and click **Generate a new API Key**
5. Click the **Token** hyperlink next to your API key, approve the permissions, and copy the token

See [Getting Started](docs/getting-started.md) for the full walkthrough.

## Quick Start

Store your credentials:

```bash
trello auth set --api-key <your-api-key> --token <your-token>
trello auth status --pretty
```

Or use interactive browser login:

```bash
trello auth set-key --api-key <your-api-key>
trello auth login
```

List your boards:

```bash
trello boards list --pretty
```

## Output Shape

Success:

```json
{"ok":true,"data":{"version":"1.0.0","commit":"abc1234","date":"2026-03-13"}}
```

Error:

```json
{"ok":false,"error":{"code":"VALIDATION_ERROR","message":"--board is required"}}
```

Use `--pretty` on any command to indent the JSON.

## Common Commands

```bash
trello boards list
trello boards create --name "Project Alpha" --default-lists
trello lists list --board <board-id>
trello cards list --list <list-id>
trello cards create --list <list-id> --name "Follow up with customer"
trello search cards --query "customer"
trello custom-fields list --board <board-id>
trello custom-fields items set --card <card-id> --field <field-id> --text "value"
```

## Documentation

- [Getting Started](docs/getting-started.md)
- [Authentication](docs/concepts/authentication.md)
- [JSON Output](docs/concepts/json-output.md)
- [Errors](docs/concepts/errors.md)
- [Configuration](docs/concepts/configuration.md)
- [Command Reference](docs/commands/README.md)
- [Examples](docs/examples/create-a-card.md)
- [LLM Digest](LLM.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, PR process, and guidelines.

## License

[MIT](LICENSE)
