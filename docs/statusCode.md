# Status codes

SMS.ir returns these in the JSON envelope field `status`. On failure the client returns `*APIError`: `Status` is this code, `Message` is the provider text, and `HTTPStatus` is the HTTP code.

A successful call has `status == 1` (`StatusOK`).

Codes with no constant are still returned as `APIError.Status`. Compare the number directly.

## Request status

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `StatusError` | The request failed |
| 1 | `StatusOK` | Success |
| 10 | `StatusInvalidAPIKey` | Invalid API key |
| 11 | `StatusInactiveAPIKey` | Inactive API key |
| 12 | `StatusIPRestricted` | API key is limited to configured IPs |
| 13 | `StatusAccountDisabled` | Account is disabled |
| 14 | `StatusAccountSuspended` | Account is suspended |
| 15 | `StatusUpgradePlan` | Upgrade the plan to use the web service |
| 16 | `StatusBadParameter` | A parameter value is invalid |
| 20 | `StatusRateLimited` | Too many requests |
| 101 | `StatusInvalidLine` | Invalid line number |
| 102 | `StatusInsufficientCredit` | Insufficient credit |
| 103 | `StatusEmptyText` | Message text is empty |
| 104 | `StatusInvalidMobile` | Invalid mobile number |
| 105 | `StatusTooManyMobiles` | More than 100 mobiles |
| 106 | | More than 100 message texts |
| 107 | | Mobile list is empty |
| 108 | | Message text list is empty |
| 109 | | Invalid send time |
| 110 | | Mobile count and text count differ |
| 111 | | No send is registered with this id |
| 112 | | Nothing found to delete |
| 113 | `StatusTemplateNotFound` | Template not found |
| 114 | `StatusParamTooLong` | Parameter value is longer than 25 characters |
| 115 | `StatusBlacklisted` | Mobile is blacklisted |
| 116 | `StatusMissingParameter` | A verify parameter name was not set |
| 117 | | Message text was not approved |
| 118 | | Too many messages |
| 119 | | Upgrade the plan to use a custom template |
| 121 | | Upgrade the plan to add more templates |
| 122 | | A template under review cannot be edited |
| 123 | | Sender line needs activation |
| 124 | | Only OTP sends are allowed, and this template is not an OTP template |
| 125 | | Template text must include at least one `#PARAMETER#` variable |

## Delivery status

`MessageReport.DeliveryState`. JSON `null` decodes as `0`, which means the state is not set yet.

| Code | Constant | Meaning |
|------|----------|---------|
| 1 | `DeliveryReachedPhone` | Delivered to the phone |
| 2 | `DeliveryNotReachedPhone` | Not delivered to the phone |
| 3 | `DeliveryReachedTelecom` | Reached telecom |
| 4 | `DeliveryNotReachedTelecom` | Not reached telecom |
| 5 | `DeliveryReachedOperator` | Reached the operator |
| 6 | `DeliveryFailed` | Failed |
| 7 | `DeliveryBlacklist` | Blacklisted |
| 8 | `DeliveryUnknown` | Unknown |
