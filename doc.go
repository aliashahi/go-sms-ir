// Package smsir is an unofficial typed Go client for the SMS.ir REST API (https://api.sms.ir).
//
// It is not affiliated with SMS.ir or IPeCompany. It covers the same surface as the official Python SDK
// (https://github.com/IPeCompany/SmsPanelV2.Python): bulk and like-to-like
// sends, verify/OTP templates (list, get, create, update, and delete), reports,
// inbound messages, credit, and lines.
//
// The package depends only on the Go standard library. Pass credentials into
// [New]; it never reads environment variables or application config.
//
// Quick start:
//
//	client := smsir.New("YOUR_API_KEY", 30004505000017)
//	res, err := client.SendVerify(ctx, smsir.VerifyRequest{
//		Mobile:     "09121234567",
//		TemplateID: 123456,
//		Parameters: []smsir.VerifyParameter{{Name: "CODE", Value: "123456"}},
//	})
package smsir
