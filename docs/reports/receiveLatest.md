# Latest unread inbound

Unread inbound messages. SMS.ir marks each returned message as read, and this method will not return it again.

`Client.LatestReceived` · `GET /v1/receive/latest`

```go
rows, err := client.LatestReceived(ctx, 20, 30004505000017)
```

`count` `0` omits the parameter. The API default and maximum are 100. `lineNumber` `0` returns unread messages on every line for this API key. A non-zero line must be active and belong to that key.

## Query

| Name | Go | Required | Notes |
|------|----|----------|-------|
| `count` | `count int` | no | Max 100, default 100 |
| `lineNumber` | `lineNumber int64` | no | `0` means every line on the key |

## Response

| JSON | Go | Notes |
|------|----|-------|
| `receiveReturnId` | `ReceiveReturnID int64` | Inbound message id |
| `messageText` | `MessageText string` | |
| `number` | `Number int64` | Receiving line |
| `mobile` | `Mobile FlexString` | Sender. JSON number or string |
| `receivedDateTime` | `ReceivedDateTime int64` | Unix time, UTC |

## Example

`GET https://api.sms.ir/v1/receive/latest?count=20&lineNumber=30004505000017`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "receiveReturnId": 123456789,
      "messageText": "HyperBox.ir فروشگاه اینترنتی اسباب بازی و عروسک",
      "number": 30004505000017,
      "mobile": 91200000002,
      "receivedDateTime": 1628683625
    }
  ]
}
```
