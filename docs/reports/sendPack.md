# Pack report

Messages in one send pack, with delivery status.

`Client.PackReport` · `GET /v1/send/pack/{packId}`

`packId` comes from `SendBulk` or `SendLikeToLike`. This is not the list of today's packs. That list is `TodayPacks`.

```go
rows, err := client.PackReport(ctx, "bdec19c9-2736-4095-8ef1-ea21afe3771f")
```

## Response

`data` is an array of [message reports](sendReports.md). `mobile` may be a JSON number. A null `deliveryState` decodes as `0`.

| JSON | Go | Notes |
|------|----|-------|
| `messageId` | `MessageID int64` | |
| `mobile` | `Mobile FlexString` | JSON number or string |
| `messageText` | `MessageText string` | |
| `sendDateTime` | `SendDateTime int64` | Unix time, UTC |
| `lineNumber` | `LineNumber int64` | Sender line |
| `cost` | `Cost float64` | Credit charged |
| `deliveryState` | `DeliveryState int` | Nullable. JSON `null` is `0` |
| `deliveryDateTime` | `DeliveryDateTime *int64` | Nullable Unix time |

Delivery codes are in [statusCode.md](../statusCode.md).

## Example

`GET https://api.sms.ir/v1/send/pack/bdec19c9-2736-4095-8ef1-ea21afe3771f`

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
