# Custom Fields CLI Design

## Summary

Add first-class custom field support to the Trello CLI for both:

- board-level custom field definitions
- card-level custom field item values

The implementation should introduce a new top-level `custom-fields` command group, extend the Trello client with JSON-body mutation helpers for custom field endpoints, add typed request and response models, and update agent-facing and user-facing documentation.

## Goals

- Support listing, reading, creating, updating, and deleting custom field definitions on boards
- Support managing list options for list-type custom fields
- Support listing, setting, and clearing custom field values on cards
- Preserve the existing CLI conventions:
  - grouped Cobra commands by resource
  - strict local flag validation
  - stable JSON envelopes
  - small typed client methods in `internal/trello`
- Update agent-facing documentation so LLMs and skills know the new command surface

## Non-Goals

- Power-Up administration beyond custom fields
- Bulk value updates across multiple cards in one command
- Name-based resolution for fields or options
- Auto-discovery of field type before value validation in the CLI

## API Surface

The feature is based on Trello's custom fields REST APIs:

- board-level custom field definitions
- custom field option management for list-type fields
- card-level custom field items

Important API characteristic:

- custom field mutations use JSON request bodies rather than only query-string mutations

This requires a small client enhancement so the CLI can keep its current handler structure while supporting the new endpoints cleanly.

## Command Design

### Top-Level Group

Add a new top-level command:

- `trello custom-fields`

This command should follow the existing resource-group pattern already used by `boards`, `cards`, `checklists`, `labels`, and others.

### Definition Commands

- `trello custom-fields list --board <board-id>`
- `trello custom-fields get --field <custom-field-id>`
- `trello custom-fields create --board <board-id> --name <name> --type <text|number|date|checkbox|list> [--pos <top|bottom|n>] [--card-front] [--option <text> ...]`
- `trello custom-fields update --field <custom-field-id> [--name <name>] [--pos <top|bottom|n>] [--card-front=<bool>]`
- `trello custom-fields delete --field <custom-field-id>`

### Option Commands

- `trello custom-fields options list --field <custom-field-id>`
- `trello custom-fields options add --field <custom-field-id> --text <label> [--color <color>] [--pos <top|bottom|n>]`
- `trello custom-fields options update --field <custom-field-id> --option <option-id> [--text <label>] [--color <color>] [--pos <top|bottom|n>]`
- `trello custom-fields options delete --field <custom-field-id> --option <option-id>`

### Card Item Commands

- `trello custom-fields items list --card <card-id>`
- `trello custom-fields items set --card <card-id> --field <custom-field-id>` with exactly one of:
  - `--text <value>`
  - `--number <value>`
  - `--date <iso-8601>`
  - `--checked <true|false>`
  - `--option <option-id>`
- `trello custom-fields items clear --card <card-id> --field <custom-field-id>`

## Why This Shape

Three shapes were considered:

1. A dedicated `custom-fields` command group
2. Splitting board definition commands under `boards` and card value commands under `cards`
3. Shipping a reduced read/write subset first

The dedicated `custom-fields` group is preferred because it:

- matches the Trello resource concept closely enough for users
- keeps board definitions and card values discoverable in one place
- avoids scattering tests and docs across unrelated command trees
- aligns with the repository's current grouped-command style

## Validation Rules

Validation should remain in the Cobra layer, following current project practice.

### General

- `list` requires `--board`
- `get` requires `--field`
- `delete` requires `--field`
- item `list` requires `--card`
- item `clear` requires both `--card` and `--field`

### Create

- requires `--board`, `--name`, and `--type`
- `--type` must be one of `text`, `number`, `date`, `checkbox`, or `list`
- repeated `--option` values are allowed only when `--type list`
- non-list field creation must reject `--option`

### Update

- requires `--field`
- requires at least one mutation flag

### Options

- option `list` requires `--field`
- option `add` requires `--field` and `--text`
- option `update` requires `--field`, `--option`, and at least one mutation flag
- option `delete` requires `--field` and `--option`

### Item Values

- item `set` requires `--card` and `--field`
- item `set` must accept exactly one value mode among `--text`, `--number`, `--date`, `--checked`, or `--option`
- `--date` must use the existing ISO-8601 validation rules
- `--checked` must parse as a boolean

## Data Model

Add the following Trello types to `internal/trello/types.go`:

- `CustomField`
- `CustomFieldDisplay`
- `CustomFieldOption`
- `CustomFieldOptionValue`
- `CustomFieldItem`
- `CustomFieldItemValue`

Add request structs for:

- custom field creation
- custom field update
- option creation
- option update
- item set

The card item set request should model Trello's type-specific payload shape explicitly rather than hiding it behind untyped maps. This will make command validation and tests easier to understand.

## Client Changes

Add a new file:

- `internal/trello/custom_fields.go`

Add API interface methods in `internal/trello/client.go` for:

- listing board custom fields
- getting a custom field
- creating a custom field
- updating a custom field
- deleting a custom field
- listing field options
- creating an option
- updating an option
- deleting an option
- listing card custom field items
- setting a card custom field item
- clearing a card custom field item

### JSON Request Helpers

Add JSON-body mutation helpers to `internal/trello/client.go`:

- `PostJSON`
- `PutJSON`

These helpers should:

- send the same authenticated requests as existing methods
- set `Content-Type: application/json`
- decode JSON responses into typed structs
- preserve the current retry and error-mapping behavior

Existing query-string methods should remain unchanged for the rest of the CLI.

## CLI Structure

Add a new command file:

- `cmd/trello/custom_fields.go`

This file should define:

- `customFieldsCmd`
- definition subcommands
- `options` subgroup
- `items` subgroup

Command handlers should:

- resolve auth the same way as existing commands
- validate flags before client calls
- construct typed parameter structs
- output the same success envelope via `output`

## Testing Strategy

### Command Tests

Add:

- `cmd/trello/custom_fields_test.go`

Cover:

- required flag validation
- mutually exclusive item value flags
- list-type option validation
- bool and ISO-8601 validation for values
- successful handler wiring through the mocked `trello.API`
- JSON envelope shape on success and failure

### Client Tests

Add:

- `internal/trello/custom_fields_test.go`

Cover:

- correct endpoints and methods
- JSON body encoding for definition, option, and item mutations
- query params where still required by Trello
- response decoding

Update:

- `internal/trello/client_test.go`

Cover:

- `PostJSON`
- `PutJSON`

## Documentation Updates

Update the following files as part of the implementation:

- `README.md`
  - add custom fields to the supported feature list
  - add a few representative examples
- `LLM.md`
  - add the `custom-fields` command group
  - document item set value modes and common workflow guidance
- `docs/commands/README.md`
  - link the new command reference page
- `docs/commands/custom-fields.md`
  - add full command reference, rules, and examples

The user also requested a skills documentation update. That should include the Trello CLI skill:

- `/Users/brettmcdowell/.codex/skills/using-trello-cli/SKILL.md`

The skill update should expand the command/reference guidance so future agent flows can discover and use `custom-fields` safely.

Also update the existing command-spec document to remove `custom fields` from the V1 exclusions now that the repo will support it.

## Risks And Decisions

### Value-Type Ambiguity

The CLI will not infer a field's type automatically before setting a value. Users must choose the matching value flag. This keeps the initial implementation deterministic and testable.

### List Option Lifecycle

List-type fields are only practical if options can be managed, so option commands are included in the first implementation rather than deferred.

### Backward Compatibility

This change is additive. No existing command shape needs to change.

## Implementation Notes

- Prefer repeated `--option` flags for list-field creation over comma-separated parsing
- Keep boolean parsing explicit for `--checked` and `--card-front`
- Reuse the existing contract helpers where possible
- Follow existing file naming patterns using underscores in Go file names only where already accepted by the repo style

## Deliverables

- new `custom-fields` CLI command group
- client support for Trello custom field endpoints
- JSON-body request support in the Trello client
- tests for command and client layers
- updated command docs, README, and LLM digest
- updated Trello CLI skill documentation
