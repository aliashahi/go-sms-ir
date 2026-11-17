# Archived sent messages

Sent messages through the end of the previous day. Today's sends are on [sendLive.md](sendLive.md).

`Client.ArchiveReport` · `GET /v1/send/archive`

```go
from := int64(1613465574)
to := int64(1623805200)
rows, err := client.ArchiveReport(ctx, &from, &to, 100, 1)
```

Pass a nil date to omit it. Pass `0` for page size or page number to omit those too. `pageSize` cannot exceed 100. Dates are Unix time, UTC.

## Query

| Name | Go | Required | Notes |
|------|----|----------|-------|
| `fromDate` | `fromDate *int64` | no | Start, inclusive |
| `toDate` | `toDate *int64` | no | End |
| `pageSize` | `pageSize int` | no | Max 100, default 100 |
| `pageNumber` | `pageNumber int` | no | Default 1 |

## Response

`data` is an array of message reports. Fields match [sendPack.md](sendPack.md).

## Example

`GET https://api.sms.ir/v1/send/archive?pageNumber=1&pageSize=100&fromDate=1613465574&toDate=1623805200`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "messageId": 89545112,
      "mobile": 91200000677,
      "messageText": "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
      "sendDateTime": 1628583626,
      "lineNumber": 30004505000017,
      "cost": 1.0,
      "deliveryState": 1,
      "deliveryDateTime": 1628683629
    }
  ]
}
```
