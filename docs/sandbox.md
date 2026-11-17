# Sandbox

The sandbox is the same API with a sandbox key. Requests are validated like production, but no SMS is sent, no credit is charged, and nothing is stored in reports. Returned data is simulated.

Create the key in the panel: developers, API keys, new key, type Sandbox. Pass it to `New`. URLs, bodies, and errors stay the same.

```go
client := smsir.New(sandboxAPIKey, 0)
res, err := client.SendVerify(ctx, smsir.VerifyRequest{
	Mobile:     "09190000004",
	TemplateID: smsir.SandboxVerifyTemplateID,
	Parameters: []smsir.VerifyParameter{{Name: "CODE", Value: "12345"}},
})
```

## Built-in verify template

| | |
|--|--|
| ID | `123456` (`SandboxVerifyTemplateID`) |
| Text | `کد تایید شما: #CODE#` |

Use this when the site is not ready for a real template. `SendVerify` with the sandbox key accepts this id.

## Templates

`Templates`, `Template`, `AddTemplate`, `ModifyTemplate`, and `RemoveTemplate` do not create, edit, or delete templates in the real panel. List and get return the sandbox templates only. Add returns the sandbox template id, which `SendVerify` can use.

## What is not stored

Send and report calls return a response and stop there. The panel report pages stay empty.
