# Archived inbound messages

Inbound messages through the end of the previous day. Today's inbound messages are on [receiveLive.md](receiveLive.md).

`Client.ArchiveReceived` · `GET /v1/receive/archive`

```go
from := int64(1613465574)
to := int64(1623805200)
rows, err := client.ArchiveReceived(ctx, &from, &to, 100, 1, "09120000002")
```

Pass a nil date to omit it. Pass `0` for page size or page number to omit those too. `pageSize` cannot exceed 100. Pass `""` for `mobile` to list every sender. Dates are Unix time, UTC.

## Query

| Name | Go | Required | Notes |
|------|----|----------|-------|
| `fromDate` | `fromDate *int64` | no | Start |
| `toDate` | `toDate *int64` | no | End |
| `pageSize` | `pageSize int` | no | Max 100, default 100 |
| `pageNumber` | `pageNumber int` | no | Default 1 |
| `mobile` | `mobile string` | no | Sender mobile |

## Response

Same fields as [receiveLatest.md](receiveLatest.md): `receiveReturnId`, `messageText`, `number`, `mobile`, `receivedDateTime`.

## Example

`GET https://api.sms.ir/v1/receive/archive?pageNumber=1&pageSize=100&fromDate=1613465574&toDate=1623805200`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "receiveReturnId": 12345678987,
      "messageText": "HyperBox.ir فروشگاه اینترنتی اسباب بازی و عروسک",
      "number": 30004505000017,
      "mobile": 91200000002,
      "receivedDateTime": 1628683625
    }
  ]
}
```
