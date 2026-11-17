# Verify send

High-priority templated SMS, including OTP and invoices. It goes out on a service line, so it can reach numbers that blocked advertising SMS.

`Client.SendVerify` · `POST /v1/send/verify`

The template must already be approved. Create one with `AddTemplate` or in the panel, then wait for `TemplateStatusApproved`. See [templates.md](../template/templates.md).

```go
res, err := client.SendVerify(ctx, smsir.VerifyRequest{
	Mobile:     "09190000004",
	TemplateID: 123456,
	Parameters: []smsir.VerifyParameter{{Name: "CODE", Value: "12345"}},
})
```

A line number is not used. `Name` is the placeholder without `#`. `Value` is at most 25 characters. The client strips a surrounding `#` from `Name`.

In the sandbox, the built-in template is `SandboxVerifyTemplateID` (`123456`), text `کد تایید شما: #CODE#`. See [sandbox.md](../sandbox.md).

## Body

| JSON | Go | Required | Notes |
|------|----|----------|-------|
| `mobile` | `Mobile string` | yes | Normalized with `FormatMobile` |
| `templateId` | `TemplateID int` | yes | Approved template id |
| `parameters` | `Parameters []VerifyParameter` | yes | At least one |

`VerifyParameter` is `name` (`Name`) and `value` (`Value`, max 25 characters).

## Response

| JSON | Go | Notes |
|------|----|-------|
| `messageId` | `MessageID int64` | |
| `cost` | `Cost float64` | Credit charged |

## Example

```json
{
  "mobile": "09190000004",
  "templateId": 123456,
  "parameters": [
    {"name": "CODE", "value": "12345"}
  ]
}
```

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
