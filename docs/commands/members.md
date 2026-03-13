# `trello members`

Manage member lookups on boards and member assignment on cards.

## Subcommands

- `trello members list --board <board-id>`
- `trello members add --card <card-id> --member <member-id>`
- `trello members remove --card <card-id> --member <member-id>`

## Examples

```bash
trello members list --board <board-id>
trello members add --card <card-id> --member <member-id>
trello members remove --card <card-id> --member <member-id>
```

## Usage Pattern

Use `members list` to discover valid member IDs before assigning them to cards.
