# go-sms-ir

[![Go Reference](https://pkg.go.dev/badge/github.com/aliashahi/go-sms-ir.svg)](https://pkg.go.dev/github.com/aliashahi/go-sms-ir)
[![Go Report Card](https://goreportcard.com/badge/github.com/aliashahi/go-sms-ir)](https://goreportcard.com/report/github.com/aliashahi/go-sms-ir)

Unofficial typed Go client for the [SMS.ir](https://sms.ir) REST API (`https://api.sms.ir`).

This project is an independent library. It is not affiliated with, endorsed by, or maintained by SMS.ir or IPeCompany. Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).

- Standard library only (no third-party deps)
- Context-aware HTTP calls
- Bulk, like-to-like, verify/OTP, reports, inbound, credit, lines, templates
- IR mobile normalization (`+98…` / `912…` / Persian digits → `09…`)
- Structured `APIError` with HTTP + business status codes

Compatible with the official Python SDK surface: [IPeCompany/SmsPanelV2.Python](https://github.com/IPeCompany/SmsPanelV2.Python).

## Install

```bash
go get github.com/aliashahi/go-sms-ir
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	smsir "github.com/aliashahi/go-sms-ir"
)

func main() {
	ctx := context.Background()
	client := smsir.New(
		os.Getenv("SMSIR_API_KEY"),
		0, // line number only needed for bulk / like-to-like
	)

	res, err := client.SendVerify(ctx, smsir.VerifyRequest{
		Mobile:     "09121234567",
		TemplateID: 123456, // approved template id from the panel
		Parameters: []smsir.VerifyParameter{
			{Name: "CODE", Value: "123456"},
		},
	})
	if err != nil {
		if apiErr, ok := smsir.AsAPIError(err); ok {
			log.Fatalf("sms.ir status=%d http=%d: %s", apiErr.Status, apiErr.HTTPStatus, apiErr.Message)
		}
		log.Fatal(err)
	}

	fmt.Println("message id:", res.MessageID, "cost:", res.Cost, "http:", res.HTTPStatus)
}
```

The import path is `github.com/aliashahi/go-sms-ir`; the package name is `smsir`.

## Configuration

| Input | Required | Notes |
|-------|----------|--------|
| API key | yes | Panel → developers → `X-API-KEY` |
| Line number | for bulk / like-to-like | Ignored by `SendVerify` |
| Template id | for verify | Must be **approved** in panel fast-send |

Options:

```go
client := smsir.New(
	apiKey,
	lineNumber,
	smsir.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	smsir.WithTimeout(10 * time.Second), // sets default client timeout
	smsir.WithBaseURL("https://api.sms.ir"),
	smsir.WithUserAgent("my-app/1.0"),
)
```

`client.Configured()` is `false` when the API key is empty — useful for optional integrations / fallbacks.

## API map

| Method | HTTP | Description |
|--------|------|-------------|
| `SendSMS` | `POST /v1/send/bulk` | Single free-text SMS (wrapper around bulk) |
| `SendBulk` | `POST /v1/send/bulk` | One text → many mobiles (max 100) |
| `SendLikeToLike` | `POST /v1/send/likeToLike` | Pair messages ↔ mobiles |
| `SendVerify` | `POST /v1/send/verify` | Templated / OTP (high priority) |
| `CancelScheduled` | `DELETE /v1/send/scheduled/{packId}` | Cancel scheduled pack |
| `MessageReport` | `GET /v1/send/{messageId}` | One message report |
| `PackReport` | `GET /v1/send/pack/{packId}` | Pack report |
| `TodayReport` | `GET /v1/send/live` | Today’s outbound |
| `ArchiveReport` | `GET /v1/send/archive` | Archived outbound |
| `LatestReceived` | `GET /v1/receive/latest` | Latest inbound |
| `TodayReceived` | `GET /v1/receive/live` | Today’s inbound |
| `ArchiveReceived` | `GET /v1/receive/archive` | Archived inbound |
| `Credit` | `GET /v1/credit` | Account credit |
| `Lines` | `GET /v1/line` | Sender lines |
| `Templates` | `GET /v1/templates` | Verify templates |

Official docs: [sms.ir REST API](https://sms.ir/rest-api/).

## Verify / OTP

1. Create a template in the panel (fast-send), e.g. `کد تایید شما #CODE# می‌باشد.`
2. Wait until status is **approved** (`TemplateStatusApproved = 2`).
3. Call `SendVerify` with `TemplateID` and parameters (`Name` without `#`).

```go
approved := smsir.TemplateStatusApproved
templates, err := client.Templates(ctx, 1, 50, &approved)
```

Parameter values are limited by SMS.ir (typically ≤ 25 characters).

## Bulk send

```go
res, err := client.SendBulk(ctx, smsir.BulkSendRequest{
	LineNumber:  30004505000017,
	MessageText: "Hello from go-sms-ir",
	Mobiles:     []string{"+989121234567", "09121234568"},
})
```

Notes from the API:

- Max **100** mobiles per request
- `MessageIDs[i] == 0` → that mobile is blacklisted
- Optional `SendDateTime` is Unix time (UTC)

## Phone numbers

`FormatMobile` (also applied automatically on send) accepts:

| Input | Output |
|-------|--------|
| `+989121234567` | `09121234567` |
| `989121234567` | `09121234567` |
| `9121234567` | `09121234567` |
| `۰۹۱۲۱۲۳۴۵۶۷` | `09121234567` |

## Errors

Successful responses have envelope `status == 1`.

Failures return `*smsir.APIError`:

```go
res, err := client.SendVerify(ctx, req)
if err != nil {
	if apiErr, ok := smsir.AsAPIError(err); ok {
		// apiErr.HTTPStatus — HTTP code
		// apiErr.Status     — SMS.ir business code
		// apiErr.Message    — provider message
		_ = apiErr
	}
	return err
}
_ = res
```

Common business codes (partial):

| Code | Meaning |
|------|---------|
| 0 | Error |
| 1 | OK |
| 10 | Invalid API key |
| 11 | Inactive API key |
| 12 | IP restricted |
| 20 | Rate limited |
| 102 | Insufficient credit |
| 104 | Invalid mobile |
| 113 | Template not found |
| 114 | Parameter value too long |
| 115 | Blacklisted |
| 116 | Missing parameter |

Delivery states on reports: `DeliveryReachedPhone` (1), …, `DeliveryBlacklist` (7), `DeliveryUnknown` (8).

Empty API key → `smsir.ErrNotConfigured`.

## Optional integration pattern

Keep your existing OTP mechanism when credentials are missing:

```go
if !client.Configured() || templateID == 0 {
	// fallback: local/dev OTP, email, etc.
	return nil
}
_, err := client.SendVerify(ctx, req)
```

For production apps, prefer enqueueing sends on a worker channel so login HTTP handlers stay fast and never expose provider errors to end users.

## Testing

```bash
go test ./...
```

Tests use `httptest` and do not call the live API.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)

## Links

- [SMS.ir](https://sms.ir)
- [REST API docs](https://sms.ir/rest-api/)
- [Official Python SDK](https://github.com/IPeCompany/SmsPanelV2.Python)
- [Repository](https://github.com/aliashahi/go-sms-ir)
