package smsir

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FlexString unmarshals a JSON string or number into a Go string.
// SMS.ir sometimes returns mobile/line values as numbers (e.g. 98912…).
type FlexString string

func (s FlexString) String() string { return string(s) }

func (s *FlexString) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*s = ""

		return nil
	}
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		*s = FlexString(str)

		return nil
	}
	var num json.Number
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("smsir: flex string: %w", err)
	}
	*s = FlexString(num.String())

	return nil
}

func (s FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}
