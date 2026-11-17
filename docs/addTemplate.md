# Add a template

Register a verify template without using the panel. The text must contain at least one `#PARAMETER#` placeholder. SMS.ir reads those names from the text. `parameters` only supplies descriptions, and names that are not in the text are ignored.

A new template is pending review. `SendVerify` accepts it only after approval. How many templates you can register depends on the plan.

`Client.AddTemplate` · `POST /v1/templates`

```go
id, err := client.AddTemplate(ctx, smsir.TemplateRequest{
	Title:        "کد تایید ورود",
	TemplateText: "کد تایید شما #CODE# می‌باشد.",
	Type:         smsir.TemplateTypeOTP,
	Parameters:   []smsir.TemplateParameter{{Name: "CODE", Description: "کد یکبار مصرف"}},
})
```

`id` is the template id to pass to `SendVerify` after approval.

## Body

| JSON | Go | Required | Notes |
|------|----|----------|-------|
| `title` | `Title` | yes | Max 32 characters |
| `templateText` | `TemplateText` | yes | At least one `#PARAMETER#`. Example: `کد تایید شما #CODE# می‌باشد.` |
| `type` | `Type` | no | `0` omits the field. SMS.ir then uses `100` (`TemplateTypeOther`) |
| `parameters` | `Parameters` | no | Descriptions only. An empty list is omitted |

`TemplateParameter.Name` has no surrounding `#` and is at most 40 characters. `Description` is optional.

Template types: [templates.md](templates.md).

## Response

`data` is the new template id (`int`).

## Example

```json
{
  "title": "کد تایید ورود",
  "templateText": "کد تایید شما #CODE# می‌باشد.",
  "type": 1,
  "parameters": [
    {
      "name": "CODE",
      "description": "کد یکبار مصرف"
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
