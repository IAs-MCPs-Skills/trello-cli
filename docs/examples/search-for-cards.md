# Search For Cards

## Goal

Find cards by free-text query when you do not have an ID yet.

## Steps

Search:

```bash
trello search cards --query "documentation"
```

Inspect a matching card:

```bash
trello cards get --card <card-id> --pretty
```

Continue with related operations:

```bash
trello comments list --card <card-id>
trello attachments list --card <card-id>
```
