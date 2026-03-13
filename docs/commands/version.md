# `trello version`

Print build metadata in the standard success envelope.

## Example

```bash
trello version --pretty
```

## Response

```json
{
  "ok": true,
  "data": {
    "version": "1.0.0",
    "commit": "abc1234",
    "date": "2026-03-13T00:00:00Z"
  }
}
```

Version, commit, and date are injected at build time. Local development builds show `"dev"`, `"unknown"`, `"unknown"`. Release builds (via goreleaser) show the git tag, commit SHA, and build timestamp.
