# Bulk send

Send one text to many mobiles. Set `SendDateTime` to schedule the pack. Leave it nil to send now.

`Client.SendBulk` · `POST /v1/send/bulk`

`SendSMS` is the same call with one mobile.

```go
res, err := client.SendBulk(ctx, smsir.BulkSendRequest{
	LineNumber:  30004505000017,
	MessageText: "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
	Mobiles:     []string{"09120000001", "09190000002"},
})
```

At most 100 mobiles. A scheduled time cannot be in the past and cannot be more than 365 days ahead. Cancel a scheduled pack with [deleteScheduled.md](deleteScheduled.md).

Mobiles are passed through `FormatMobile` (`+98912…`, `912…`, and Persian digits become `09…`).

## Body

| JSON | Go | Required | Notes |
|------|----|----------|-------|
| `lineNumber` | `LineNumber int64` | yes | Falls back to the line set on `New` |
| `messageText` | `MessageText string` | yes | |
| `mobiles` | `Mobiles []string` | yes | Max 100 |
| `sendDateTime` | `SendDateTime *int64` | no | Unix time, UTC. Nil sends now |

## Response

| JSON | Go | Notes |
|------|----|-------|
| `packId` | `PackID string` | Pack id |
| `messageIds` | `MessageIDs []int64` | One id per mobile. `0` means that mobile is blacklisted |
| `cost` | `Cost float64` | Credit charged |

## Example

```json
{
  "lineNumber": 30004505000017,
  "messageText": "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
  "mobiles": ["09120000001", "09190000002"]
}
```

```json
{
  "status": 1,
  "message": "موفق",
  "data": {
    "packId": "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1",
    "messageIds": [86522023, 86522024],
    "cost": 2
  }
}
```
