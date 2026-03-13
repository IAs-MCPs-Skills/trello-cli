# `trello cards`

Manage Trello cards.

## Subcommands

- `trello cards list --board <board-id>`
- `trello cards list --list <list-id>`
- `trello cards get --card <card-id>`
- `trello cards create --list <list-id> --name <name> [--desc <text>] [--due <iso-8601>]`
- `trello cards update --card <card-id> [--name <name>] [--desc <text>] [--due <iso-8601>] [--labels <csv>] [--members <csv>]`
- `trello cards move --card <card-id> --list <list-id> [--pos <number>]`
- `trello cards archive --card <card-id>`
- `trello cards delete --card <card-id>`

## Rules

- `list` requires exactly one of `--board` or `--list`
- `create` requires `--list` and `--name`
- `update` requires `--card` and at least one mutation flag
- `--due` must be a valid ISO-8601 timestamp or date

## Examples

```bash
trello cards list --list <list-id> --pretty
trello cards create --list <list-id> --name "Write docs" --desc "First pass"
trello cards update --card <card-id> --due 2026-03-20T17:00:00Z
trello cards move --card <card-id> --list <done-list-id> --pos 1
trello cards archive --card <card-id>
trello cards delete --card <card-id>
```

## Usage Pattern

For most workflows: discover board -> discover list -> create or move card -> optionally add comments, labels, members, checklists, or attachments.
