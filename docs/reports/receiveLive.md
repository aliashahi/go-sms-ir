# Today's inbound messages

Inbound messages for today, read and unread. Early in the day this can still include yesterday.

`Client.TodayReceived` · `GET /v1/receive/live`

```go
newest := true
rows, err := client.TodayReceived(ctx, 100, 1, &newest, "09120000002")
```

Pass `0` for page size or page number to omit it. `pageSize` cannot exceed 100. Pass `nil` for `sortByNewest` to keep the API default, oldest first. Pass `""` for `mobile` to list every sender. A mobile is normalized with `FormatMobile`.

## Query

| Name | Go | Required | Notes |
|------|----|----------|-------|
| `pageSize` | `pageSize int` | no | Max 100, default 100 |
| `pageNumber` | `pageNumber int` | no | Default 1 |
| `sortByNewest` | `sortByNewest *bool` | no | `true` is newest first. Nil omits the flag |
| `mobile` | `mobile string` | no | Sender mobile |

## Response

`data` is an array of `ReceivedMessage`. This list does not include `receiveReturnId`.

| JSON | Go | Notes |
|------|----|-------|
| `messageText` | `MessageText string` | |
| `number` | `Number int64` | Receiving line |
| `mobile` | `Mobile FlexString` | Sender. JSON number or string |
| `receivedDateTime` | `ReceivedDateTime int64` | Unix time, UTC |

## Example

`GET https://api.sms.ir/v1/receive/live?pageSize=20&pageNumber=3&sortByNewest=false`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "messageText": "HyperBox.ir فروشگاه اینترنتی اسباب بازی و عروسک",
      "number": 30004505000017,
      "mobile": 91200000002,
      "receivedDateTime": 1628683625
    }
  ]
}
```
