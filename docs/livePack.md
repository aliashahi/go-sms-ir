# Today’s send packs

Summary of send packs created today. This is not the per-message report. Use `PackReport` for messages inside one pack (`GET /v1/send/pack/{packId}`).

`Client.TodayPacks` · `GET /v1/send/pack`

```go
packs, err := client.TodayPacks(ctx, 100, 1)
```

Pass `0` for either argument to omit it. SMS.ir then uses `pageSize=100` and `pageNumber=1`. `pageSize` cannot exceed 100.

## Query

| Name | Required | Notes |
|------|----------|-------|
| `pageSize` | no | Items per page. Max 100, default 100 |
| `pageNumber` | no | Page number. Default 1 |

## Response

`data` is an array.

| JSON | Go | Notes |
|------|----|-------|
| `packId` | `PackID string` | Pack id |
| `recipientCount` | `RecipientCount int` | Recipient count |
| `creationDateTime` | `CreationDateTime int64` | Unix time |

## Example

`GET https://api.sms.ir/v1/send/pack?pageNumber=1&pageSize=100`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
      "packId": "e7c09e23f0db4834b9bcb79e7b054f4c",
      "recipientCount": 100,
      "creationDateTime": 1628683626
    },
    {
      "packId": "0cf3017fd9d84babbb5ed5579104dab2",
      "recipientCount": 200,
      "creationDateTime": 1628683626
    }
  ]
}
```
