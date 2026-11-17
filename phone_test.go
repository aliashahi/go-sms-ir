package smsir

import "testing"

func TestFormatMobile(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"+989121234567", "09121234567"},
		{"989121234567", "09121234567"},
		{"09121234567", "09121234567"},
		{"9121234567", "09121234567"},
		{"۰۹۱۲۱۲۳۴۵۶۷", "09121234567"},
		{"٩١٢١٢٣٤٥٦٧", "09121234567"},
		{"", ""},
		{"  +98 912 123 4567 ", "09121234567"},
	}
	for _, tc := range cases {
		if got := FormatMobile(tc.in); got != tc.want {
			t.Fatalf("FormatMobile(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}
