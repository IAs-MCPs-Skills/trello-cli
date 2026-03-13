# Errors

## Error Shape

All command failures are returned in the standard error envelope:

```json
{
  "ok": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "--card is required"
  }
}
```

The CLI exits with a non-zero status on error.

## Error Codes

- `AUTH_REQUIRED`: credentials are missing or incomplete
- `AUTH_INVALID`: Trello rejected the credentials
- `NOT_FOUND`: requested resource was not found
- `VALIDATION_ERROR`: missing or invalid input flags
- `CONFLICT`: operation conflicts with current state
- `RATE_LIMITED`: Trello rate limiting or retry exhaustion
- `HTTP_ERROR`: transport or non-auth HTTP failure
- `FILE_NOT_FOUND`: required local file does not exist
- `UNSUPPORTED`: requested behavior is not supported
- `UNKNOWN_ERROR`: uncategorized failure

## Common Validation Cases

- Missing required ID flags such as `--board`, `--list`, or `--card`
- Supplying both `--board` and `--list` to `cards list`
- Omitting all mutation fields on update commands
- Invalid ISO-8601 values for `--due`
- Invalid URLs for `attachments add-url`
- Missing local files for `attachments add-file`
- Invalid checklist item state values outside `complete` or `incomplete`

## Recommended Handling

- Retry only when the failure is transient or external.
- Treat `VALIDATION_ERROR` as a caller fix, not a retry candidate.
- Use `auth status` or `auth login` to resolve `AUTH_REQUIRED` and `AUTH_INVALID`.
- For file uploads, validate the file path before invoking the command in larger workflows.
