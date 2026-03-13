# JSON Output

## Contract

Every command writes JSON to `stdout`.

Success envelope:

```json
{
  "ok": true,
  "data": {}
}
```

Error envelope:

```json
{
  "ok": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "--board is required"
  }
}
```

## Formatting

- Default output is compact JSON with a trailing newline.
- `--pretty` indents the JSON and still appends a trailing newline.
- Diagnostic logs from `--verbose` go to `stderr`, not `stdout`.

## Practical Guidance

- Parse `ok` first.
- On success, inspect `data`.
- On failure, inspect `error.code` and `error.message`.
- Do not rely on English help text when automation can rely on the JSON envelope.

## Typical Result Shapes

- Lists usually return arrays in `data`
- Resource fetches usually return one object in `data`
- Delete operations return confirmation objects
- Action-style commands may return `{ "success": true, "id": "..." }`

## Example

```bash
trello version --pretty
```

```json
{
  "ok": true,
  "data": {
    "version": "dev",
    "commit": "unknown",
    "date": "unknown"
  }
}
```
