# List templates

Paged list of fast-send templates. `Id` is the `TemplateID` you pass to `SendVerify`. Only an approved template can be used to send.

`Client.Templates` · `GET /v1/templates`

```go
approved := smsir.TemplateStatusApproved
templates, err := client.Templates(ctx, 1, 10, &approved)
```

Pass `0` for page number or page size to omit that query parameter. SMS.ir then uses page 1 and 10 items. `pageSize` cannot exceed 100. Pass `nil` for status to skip the filter.

## Query

| Name | Required | Notes |
|------|----------|-------|
| `pageNumber` | no | Default 1 |
| `pageSize` | no | Default 10, max 100 |
| `status` | no | Review status. See the table below |

## Response

`data` is an array of templates.

| JSON | Go | Notes |
|------|----|-------|
| `id` | `ID int` | `SendVerify` template id |
| `title` | `Title string` | |
| `templateText` | `TemplateText string` | Placeholders look like `#CODE#` |
| `type` | `Type byte` | Template type |
| `parameters` | `Parameters []TemplateParameter` | Names have no `#` |
| `status` | `Status byte` | Review status |
| `rejectionReason` | `RejectionReason *string` | Set only when rejected. JSON `null` is `nil` |
| `creationDateTime` | `CreationDateTime int64` | Unix time, UTC |

`TemplateParameter` is `name` (`Name`) and `description` (`Description`).

## Review status

| Code | Constant | Meaning |
|------|----------|---------|
| 1 | `TemplateStatusPending` | Waiting for review. It cannot be edited |
| 2 | `TemplateStatusApproved` | Approved. Usable with `SendVerify` |
| 3 | `TemplateStatusRejected` | Rejected. See `RejectionReason` |

## Template type

| Code | Constant | Meaning |
|------|----------|---------|
| 1 | `TemplateTypeOTP` | One-time password |
| 2 | `TemplateTypeInfo` | Notices and reminders |
| 3 | `TemplateTypeOrderStatus` | Order updates |
| 100 | `TemplateTypeOther` | Anything else. Default when create omits `type` |

## Example

`GET https://api.sms.ir/v1/templates?pageNumber=1&pageSize=10&status=2`

```json
{
  "status": 1,
  "message": "موفق",
  "data": [
    {
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
      "status": 2,
      "rejectionReason": null,
      "creationDateTime": 1628683626
    }
  ]
}
```
