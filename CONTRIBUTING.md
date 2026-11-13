# Contributing

Contributions are welcome. Bug reports, docs, tests, and pull requests all help.

**go-sms-ir is an unofficial library.** It is an independent Go client for the [SMS.ir](https://sms.ir) REST API. It is not affiliated with, endorsed by, or maintained by SMS.ir or IPeCompany. API behavior follows the public docs at [sms.ir/rest-api](https://sms.ir/rest-api/).

## What to work on

- API mismatches against the current SMS.ir REST docs
- Clearer examples and docs
- Tests for edge cases (phone formats, error envelopes, pagination)
- Small, focused improvements that stay on the standard library

## Development

Go 1.21 or newer. This module has no third-party dependencies — please keep it that way.

```bash
gofmt -w .
go test ./...
```

Tests use `httptest` and must not call the live API or require an API key.

## Pull requests

1. Fork the repo and create a branch from `main`.
2. Keep the change focused on one thing.
3. Run `gofmt` and `go test ./...` before you open the pull request.
4. Describe what changed and why. Link an issue if there is one.

By contributing, you agree your work is licensed under the [MIT License](LICENSE).
