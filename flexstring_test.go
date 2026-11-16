package smsir

import (
	"encoding/json"
	"testing"
)

func TestFlexStringUnmarshal(t *testing.T) {
	t.Parallel()
	cases := []struct {
		raw  string
		want string
	}{
		{`"09121234567"`, "09121234567"},
		{`989121234567`, "989121234567"},
		{`null`, ""},
	}
	for _, tc := range cases {
		var s FlexString
		if err := json.Unmarshal([]byte(tc.raw), &s); err != nil {
			t.Fatalf("raw %s: %v", tc.raw, err)
		}
		if string(s) != tc.want {
			t.Fatalf("raw %s: got %q want %q", tc.raw, s, tc.want)
		}
	}
}

func TestMessageReportMobileNumber(t *testing.T) {
	t.Parallel()
	raw := `{
		"messageId": 823933569,
		"mobile": 989121234567,
		"messageText": "hi",
		"sendDateTime": 0,
		"lineNumber": 30004505000017,
		"cost": 1.5,
		"deliveryState": 1,
		"deliveryDateTime": null
	}`
	var report MessageReport
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatal(err)
	}
	if report.MessageID != 823933569 {
		t.Fatalf("messageId=%d", report.MessageID)
	}
	if string(report.Mobile) != "989121234567" {
		t.Fatalf("mobile=%q", report.Mobile)
	}
	if report.DeliveryState != DeliveryReachedPhone {
		t.Fatalf("deliveryState=%d", report.DeliveryState)
	}
}
