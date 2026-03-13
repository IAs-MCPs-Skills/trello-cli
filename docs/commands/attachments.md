# `trello attachments`

Manage card attachments.

## Subcommands

- `trello attachments list --card <card-id>`
- `trello attachments add-file --card <card-id> --path <local-path> [--name <display-name>]`
- `trello attachments add-url --card <card-id> --url <http-or-https-url> [--name <display-name>]`
- `trello attachments delete --card <card-id> --attachment <attachment-id>`

## Validation

- `add-file` requires a local file that exists
- `add-url` requires a valid `http` or `https` URL
- `delete` requires both the card ID and attachment ID

## Examples

```bash
trello attachments list --card <card-id>
trello attachments add-file --card <card-id> --path ./brief.pdf --name "Project brief"
trello attachments add-url --card <card-id> --url https://example.com/spec --name "Spec"
trello attachments delete --card <card-id> --attachment <attachment-id>
```
