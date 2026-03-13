# Configuration

## User-Facing Configuration Today

The current CLI exposes configuration mainly through:

- Global flags
- Authentication environment variables
- Stored credentials in the OS keyring

## Global Flags

Available on most commands:

- `--pretty`: indent JSON output
- `--verbose`: enable diagnostic logging on `stderr`

Examples:

```bash
trello boards list --pretty
trello cards create --list <list-id> --name "Ship docs" --verbose
```

## Environment Variables

Supported credential environment variables:

- `TRELLO_API_KEY`
- `TRELLO_TOKEN`

These are especially useful for CI jobs, shell sessions, and one-off automation.

```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"
trello boards list
```

## Credential Persistence

- `trello auth set` stores credentials for later commands.
- `trello auth set-key` stores only the API key.
- `trello auth clear` removes stored credentials.

## Important Scope Note

The repository contains an internal non-secret config package, but the current user-facing command behavior is centered on the flags and auth inputs documented above. If more runtime config becomes part of the public CLI surface, this document should expand with the exact supported precedence and keys.
