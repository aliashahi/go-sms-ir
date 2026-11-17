# Cancel a scheduled send

Cancel a pack that has not been sent yet. Use the `packId` returned by [bulk.md](bulk.md) or [likeToLike.md](likeToLike.md).

`Client.CancelScheduled` · `DELETE /v1/send/scheduled/{packId}`

You can cancel until 3 minutes before the scheduled time.

```go
res, err := client.CancelScheduled(ctx, "2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1")
```

## Response

| JSON | Go | Notes |
|------|----|-------|
| `returnedCreditCount` | `ReturnedCreditCount float64` | Credit returned |
| `smsCount` | `SMSCount int` | Messages cancelled |

## Example

`DELETE https://api.sms.ir/v1/send/scheduled/2b99e63c-9bf8-4a21-9bfe-3f72dc1b46f1`

```json
{
  "status": 1,
  "message": "موفق",
  "data": {
    "returnedCreditCount": 10,
    "smsCount": 5
  }
}
```
