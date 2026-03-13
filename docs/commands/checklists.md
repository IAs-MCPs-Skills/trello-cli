# `trello checklists`

Manage checklists and checklist items on cards.

## Subcommands

- `trello checklists list --card <card-id>`
- `trello checklists create --card <card-id> --name <name>`
- `trello checklists delete --checklist <checklist-id>`
- `trello checklists items add --checklist <checklist-id> --name <name>`
- `trello checklists items update --card <card-id> --item <item-id> --state <complete|incomplete>`
- `trello checklists items delete --checklist <checklist-id> --item <item-id>`

## Rules

- Checklist item state must be either `complete` or `incomplete`
- Item update requires the card ID and item ID
- Item delete requires the checklist ID and item ID

## Examples

```bash
trello checklists create --card <card-id> --name "Launch checklist"
trello checklists items add --checklist <checklist-id> --name "Write announcement"
trello checklists items update --card <card-id> --item <item-id> --state complete
trello checklists items delete --checklist <checklist-id> --item <item-id>
```
