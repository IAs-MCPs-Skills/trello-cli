# Attach A File

## Goal

Upload a local file attachment to a card.

## Steps

Confirm the file exists locally, then run:

```bash
trello attachments add-file --card <card-id> --path ./brief.pdf --name "Brief"
```

List attachments on the card:

```bash
trello attachments list --card <card-id> --pretty
```

Remove an attachment if needed:

```bash
trello attachments delete --card <card-id> --attachment <attachment-id>
```
