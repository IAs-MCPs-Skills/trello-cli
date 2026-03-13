# `trello auth`

Manage authentication and credential storage.

## Subcommands

- `trello auth set --api-key <key> --token <token>`
- `trello auth set-key --api-key <key>`
- `trello auth status`
- `trello auth login`
- `trello auth clear`

## Notes

- `set` stores credentials without validating them immediately.
- `set-key` is useful before `login` or when rotating API keys.
- `status` validates configured credentials by calling Trello.
- `login` uses a browser flow and a localhost callback on port `3007`.
- `clear` removes stored credentials from the primary store.

## Examples

```bash
trello auth set --api-key "$TRELLO_API_KEY" --token "$TRELLO_TOKEN"
trello auth set-key --api-key "$TRELLO_API_KEY"
trello auth login
trello auth status --pretty
trello auth clear
```
