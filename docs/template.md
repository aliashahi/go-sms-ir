# Get a template

One template, including its text, type, parameters, and review status. `RejectionReason` is set only when the template was rejected.

`Client.Template` · `GET /v1/templates/{templateId}`

```go
tpl, err := client.Template(ctx, 123456)
```

## Response

`data` is one template. Fields match [templates.md](templates.md): `id`, `title`, `templateText`, `type`, `parameters`, `status`, `rejectionReason`, `creationDateTime`.

Review status and template type constants are listed there.

## Example

`GET https://api.sms.ir/v1/templates/123456`

```json
{
  "status": 1,
  "message": "موفق",
  "data": {
    "id": 123456,
    "title": "کد تایید ورود",
    "templateText": "کد تایید شما #CODE# می‌باشد.",
    "type": 1,
    "parameters": [
      {
        "name": "CODE",
        "description": "کد یکبار مصرف"
      }
    ],
    "status": 3,
    "rejectionReason": "متن قالب با نوع انتخاب شده مطابقت ندارد.",
    "creationDateTime": 1628683626
  }
}
```
