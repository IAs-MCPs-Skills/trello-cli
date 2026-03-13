# `trello boards`

Manage Trello boards.

## Subcommands

- `trello boards list`
- `trello boards get --board <board-id>`

## Flags

- `--board`: board ID for `get`

## Examples

```bash
trello boards list --pretty
trello boards get --board <board-id> --pretty
```

## Usage Pattern

Use `boards list` as the discovery step before drilling into lists, cards, labels, or members.
