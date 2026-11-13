package smsir

import (
	"strings"
	"unicode"
)

// FormatMobile converts IR phone numbers to a form SMS.ir accepts (09xxxxxxxxx).
//
// Accepted inputs include E.164 (+98912…), national (0912… / 912…), and
// Persian/Arabic-Indic digits. Non-digit characters (except a leading +) are stripped.
func FormatMobile(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return phone
	}

	var b strings.Builder
	b.Grow(len(phone))
	for i, r := range phone {
		switch {
		case r == '+' && i == 0:
			b.WriteRune(r)
		case unicode.IsDigit(r):
			b.WriteRune(toASCIIDigit(r))
		}
	}
	digits := b.String()

	switch {
	case strings.HasPrefix(digits, "+98"):
		rest := strings.TrimPrefix(digits, "+98")
		if strings.HasPrefix(rest, "0") {
			return rest
		}

		return "0" + rest
	case strings.HasPrefix(digits, "98") && len(digits) >= 12:
		rest := digits[2:]
		if strings.HasPrefix(rest, "0") {
			return rest
		}

		return "0" + rest
	case strings.HasPrefix(digits, "9") && len(digits) == 10:
		return "0" + digits
	default:
		return digits
	}
}

func toASCIIDigit(r rune) rune {
	switch {
	case r >= '۰' && r <= '۹': // Persian
		return '0' + (r - '۰')
	case r >= '٠' && r <= '٩': // Arabic-Indic
		return '0' + (r - '٠')
	default:
		return r
	}
}
