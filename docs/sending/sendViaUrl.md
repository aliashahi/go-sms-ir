# Send via URL

Send one SMS by putting the credentials and the message in the query string. This call does not use the `X-API-KEY` header. `Password` is the private key from the developers panel.

`Client.SendViaURL` · `GET /v1/send`

SMS.ir also accepts `POST` with the same query. The client calls `GET`.

```go
res, err := client.SendViaURL(ctx, smsir.URLSendRequest{
	Username: "MY_USERNAME",
	Password: "MY_API_KEY",
	Line:     30004505000017,
	Mobile:   "09120000001",
	Text:     "hello",
})
```

`Line` falls back to the line set on `New`. The client API key may be empty. The mobile is normalized with `FormatMobile`.

## Query

| Name | Go | Required | Notes |
|------|----|----------|-------|
| `username` | `Username` | yes | Panel username |
| `password` | `Password` | yes | Private key |
| `line` | `Line` | yes | Sender line |
| `mobile` | `Mobile` | yes | |
| `text` | `Text` | yes | Message text |

## Response

Same fields as [verify.md](verify.md): `messageId`, `cost`.

## Example

`GET https://api.sms.ir/v1/send?username=MY_USERNAME&password=MY_API_KEY&line=30004505000017&mobile=09120000001&text=hello`

```json
{
  "status": 1,
  "message": "موفق",
  "data": {
    "messageId": 89545112,
    "cost": 1
  }
}
```
