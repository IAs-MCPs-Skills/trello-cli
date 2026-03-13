# Manage Checklists

## Goal

Create a checklist, add an item, and mark it complete.

## Steps

Create the checklist:

```bash
trello checklists create --card <card-id> --name "Release"
```

Add an item:

```bash
trello checklists items add --checklist <checklist-id> --name "Publish release notes"
```

Mark the item complete:

```bash
trello checklists items update --card <card-id> --item <item-id> --state complete
```
