# Today's sent messages

Messages sent today.

`Client.TodayReport` · `GET /v1/send/live`

```go
rows, err := client.TodayReport(ctx, 100, 1)
```

Pass `0` for either argument to omit it. SMS.ir then uses `pageSize=100` and `pageNumber=1`. `pageSize` cannot exceed 100.

## Query

| Name | Required | Notes |
|------|----------|-------|
| `pageSize` | no | Max 100, default 100 |
| `pageNumber` | no | Default 1 |

## Response

`data` is an array of message reports. Fields match [sendPack.md](sendPack.md). Delivery codes are in [statusCode.md](../statusCode.md).

## Example

`GET https://api.sms.ir/v1/send/live?pageSize=25&pageNumber=3`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "messageId": 89545112,
      "mobile": 91200000677,
      "messageText": "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
      "sendDateTime": 1628683626,
      "lineNumber": 30004505000017,
      "cost": 1.0,
      "deliveryState": 1,
      "deliveryDateTime": 1628683629
    }
  ]
}
```
