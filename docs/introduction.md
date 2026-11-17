# Introduction

The client calls `https://api.sms.ir` and speaks JSON only. Pass the API key to `New`. The client sets the headers below and reads the envelope below. It does not request XML.

## Headers

| Header | Value | Role |
|--------|-------|------|
| `Accept` | `application/json` | Response format. Set by the client |
| `Content-Type` | `application/json` | Set on requests that have a body |
| `X-API-KEY` | Panel API key | Authentication for every method except `SendViaURL` |

## HTTP status

| Code | Meaning |
|------|---------|
| 200 | Success |
| 400 | Logical error |
| 401 | Authentication failed |
| 429 | Too many requests |
| 500 | Unexpected error |

A non-success HTTP status or an envelope `status` other than `1` comes back as `*APIError`. Business codes are in [statusCode.md](statusCode.md).

## Time

Timestamps are Unix time in UTC. `SendDateTime` uses the same clock.

## Envelope

Every response has this shape. `data` is an object, an array, or a scalar, depending on the method.

```json
{
  "status": 1,
  "message": "موفق",
  "data": [30004505000027, 10002166593818]
}
```

| Field | Meaning |
|-------|---------|
| `status` | Business status. `1` is success (`StatusOK`) |
| `message` | Provider text, often Persian |
| `data` | Payload |

## Authentication

Create a private key in the panel under developers. Pass it to `New`. The client sends it as `X-API-KEY`.

`SendViaURL` is the exception: it puts the panel username and that private key in the query string. See [sending/sendViaUrl.md](sending/sendViaUrl.md).

A sandbox key uses the same header and the same URLs. See [sandbox.md](sandbox.md).
