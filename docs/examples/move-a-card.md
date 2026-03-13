# Move A Card

## Goal

Move a card into another list, optionally controlling its position.

## Steps

Inspect destination lists:

```bash
trello lists list --board <board-id> --pretty
```

Move the card:

```bash
trello cards move --card <card-id> --list <destination-list-id> --pos 1
```

Verify the result:

```bash
trello cards get --card <card-id> --pretty
```
