# `trello labels`

Manage labels.

## Subcommands

- `trello labels list --board <board-id>`
- `trello labels create --board <board-id> --name <name> --color <color>`
- `trello labels add --card <card-id> --label <label-id>`
- `trello labels remove --card <card-id> --label <label-id>`

## Notes

- Label creation happens at the board level.
- Label assignment and removal happen at the card level.

## Examples

```bash
trello labels list --board <board-id>
trello labels create --board <board-id> --name "Blocked" --color red
trello labels add --card <card-id> --label <label-id>
trello labels remove --card <card-id> --label <label-id>
```
