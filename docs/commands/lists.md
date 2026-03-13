# `trello lists`

Manage Trello lists.

## Subcommands

- `trello lists list --board <board-id>`
- `trello lists create --board <board-id> --name <name>`
- `trello lists update --list <list-id> [--name <name>] [--pos <number>]`
- `trello lists archive --list <list-id>`
- `trello lists move --list <list-id> --board <board-id> [--pos <number>]`

## Rules

- `list` requires `--board`
- `create` requires `--board` and `--name`
- `update` requires `--list` and at least one mutation flag
- `move` requires both `--list` and destination `--board`

## Examples

```bash
trello lists list --board <board-id>
trello lists create --board <board-id> --name "In Review"
trello lists update --list <list-id> --name "Done" --pos 1
trello lists move --list <list-id> --board <board-id> --pos 2
trello lists archive --list <list-id>
```
