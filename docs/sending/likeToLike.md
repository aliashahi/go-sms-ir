# Like-to-like send

Send a different text to each mobile. The mobile list and the text list must be the same length.

`Client.SendLikeToLike` · `POST /v1/send/likeToLike`

```go
res, err := client.SendLikeToLike(ctx, smsir.LikeToLikeRequest{
	LineNumber:   30004505000017,
	MessageTexts: []string{"متن اول", "متن دوم"},
	Mobiles:      []string{"09120000001", "+989190000002"},
})
```

At most 100 mobiles. A scheduled time cannot be in the past and cannot be more than 365 days ahead. Cancel a scheduled pack with [deleteScheduled.md](deleteScheduled.md).

Mobiles are passed through `FormatMobile`.

## Body

| JSON | Go | Required | Notes |
|------|----|----------|-------|
| `lineNumber` | `LineNumber int64` | yes | Falls back to the line set on `New`. Sent as a number |
| `messageTexts` | `MessageTexts []string` | yes | Same length as `Mobiles` |
| `mobiles` | `Mobiles []string` | yes | Max 100 |
| `sendDateTime` | `SendDateTime *int64` | no | Unix time, UTC. Nil sends now |

## Response

Same object as [bulk.md](bulk.md): `packId`, `messageIds`, `cost`. `messageIds[i] == 0` means that mobile is blacklisted.

## Example

```json
{
  "lineNumber": 30004505000017,
  "messageTexts": [
    "سرویس پیامکی ایده پردازان با بیش از یک دهه سابقه همراه شماست",
    "ipdemy.ir پلتفرم آموزش آنلاین، آکادمی ایده پردازان"
  ],
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
