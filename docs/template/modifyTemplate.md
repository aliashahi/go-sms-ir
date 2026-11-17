# Modify a template

Change the title, text, type, or parameter descriptions of a template. Omitted optional fields stay as they are.

A template cannot be edited while it is pending review. After a successful edit it goes back to pending, and `SendVerify` cannot use it until it is approved again.

`Client.ModifyTemplate` · `PUT /v1/templates/{templateId}`

```go
id, err := client.ModifyTemplate(ctx, 123456, smsir.TemplateRequest{
	Title:        "کد تایید ورود",
	TemplateText: "کد ورود شما #CODE# می‌باشد.",
	Type:         smsir.TemplateTypeOTP,
	Parameters:   []smsir.TemplateParameter{{Name: "CODE", Description: "کد ورود دو مرحله‌ای"}},
})
```

## Body

| JSON | Go | Required | Notes |
|------|----|----------|-------|
| `title` | `Title` | yes | Max 32 characters |
| `templateText` | `TemplateText` | yes | At least one `#PARAMETER#` |
| `type` | `Type` | no | `0` omits the field and keeps the current type |
| `parameters` | `Parameters` | no | `nil` or empty keeps the current descriptions |

`TemplateParameter.Name` has no surrounding `#` and is at most 40 characters. Names that are not in the text are ignored.

Template types: [templates.md](templates.md).

## Response

`data` is the edited template id (`int`).

## Example

`PUT https://api.sms.ir/v1/templates/123456`

```json
{
  "title": "کد تایید ورود",
  "templateText": "کد ورود شما #CODE# می‌باشد.",
  "type": 1,
  "parameters": [
    {
      "name": "CODE",
      "description": "کد ورود دو مرحله‌ای"
    }
  ]
}
```

```json
{
  "status": 1,
  "message": "موفق",
  "data": 123456
}
```
