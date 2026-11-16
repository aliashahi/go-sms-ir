# Message report

Delivery status for one sent message.

`Client.MessageReport` · `GET /v1/send/{messageId}`

The `messageId` comes back from `SendSMS`, `SendBulk`, `SendLikeToLike`, or `SendVerify`.

```go
report, err := client.MessageReport(ctx, messageID)
```

## Response

`data` is one object.

| JSON | Go | Notes |
|------|----|-------|
| `messageId` | `MessageID int64` | |
| `mobile` | `Mobile FlexString` | JSON number or string |
| `messageText` | `MessageText string` | |
| `sendDateTime` | `SendDateTime int64` | Unix time |
| `lineNumber` | `LineNumber int64` | |
| `cost` | `Cost float64` | Credit charged |
| `deliveryState` | `DeliveryState int` | Nullable. JSON `null` decodes as `0` |
| `deliveryDateTime` | `DeliveryDateTime *int64` | Nullable Unix time. `nil` when absent |

Delivery codes are in [statusCode.md](statusCode.md).

## Example

`GET https://api.sms.ir/v1/send/89545112`

```json
{
  "status": 1,
  "message": "موفق",
  "data": {
    "messageId": 89545112,
    "mobile": 91200000677,
    "messageText": "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
    "sendDateTime": 1628683626,
    "lineNumber": 30004505000017,
    "cost": 1.0,
    "deliveryState": 1,
    "deliveryDateTime": 1628683629
  }
}
```
