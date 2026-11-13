package smsir_test

import (
	"fmt"

	"github.com/aliashahi/go-sms-ir"
)

func ExampleFormatMobile() {
	fmt.Println(smsir.FormatMobile("+989121234567"))
	fmt.Println(smsir.FormatMobile("9121234567"))
	fmt.Println(smsir.FormatMobile("۰۹۱۲۱۲۳۴۵۶۷"))

	// Output:
	// 09121234567
	// 09121234567
	// 09121234567
}

func ExampleNew() {
	client := smsir.New(
		"YOUR_API_KEY",
		30004505000017, // default line for bulk sends; unused by SendVerify
		smsir.WithUserAgent("my-app/1.0"),
	)
	_ = client
}
