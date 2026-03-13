# Create A Card

## Goal

Create a new card when you know the target board but not yet the list ID.

## Steps

Discover lists on the board:

```bash
trello lists list --board <board-id> --pretty
```

Create the card:

```bash
trello cards create --list <list-id> --name "Write CLI docs" --desc "First full documentation pass"
```

Optionally confirm it:

```bash
trello cards list --list <list-id> --pretty
```
