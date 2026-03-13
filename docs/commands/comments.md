# `trello comments`

Manage card comments.

## Subcommands

- `trello comments list --card <card-id>`
- `trello comments add --card <card-id> --text <text>`
- `trello comments update --action <action-id> --text <text>`
- `trello comments delete --action <action-id>`

## Notes

- Comment updates and deletes use Trello action IDs, not card IDs.
- `list` uses the card ID to discover available comments.

## Examples

```bash
trello comments list --card <card-id>
trello comments add --card <card-id> --text "Ready for review"
trello comments update --action <action-id> --text "Ready for final review"
trello comments delete --action <action-id>
```
