# Sender lines

Line numbers that are ready to send.

`Client.Lines` · `GET /v1/line`

```go
lines, err := client.Lines(ctx)
```

`lines` is `[]int64`. Use one of these values as `LineNumber` on bulk and like-to-like sends. `SendVerify` does not use a line number.

## Example

`GET https://api.sms.ir/v1/line`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    10002155613464,
    30004505000017
  ]
}
```
