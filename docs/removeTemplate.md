# Remove a template

Delete a template by id. `SendVerify` can no longer use that id. Reports for messages already sent stay as they are.

`Client.RemoveTemplate` · `DELETE /v1/templates/{templateId}`

```go
err := client.RemoveTemplate(ctx, 123456)
```

## Response

`data` is `true` when the template was deleted. A failed delete is an `*APIError`.

## Example

`DELETE https://api.sms.ir/v1/templates/123456`

```json
{
  "status": 1,
  "message": "موفق",
  "data": true
}
```
