package request

import (
	"testing"
)

type validateTimeCase struct {
	name      string
	input     string
	wantError bool
	allowZero bool
}

var validateTimeTests = []validateTimeCase{
	{"valid date", "2025-10-31 14:30", false, false},
	{"invalid format", "25-10-2025 14:30", true, true},
	{"invalid date values", "2025-13-32 14:30", true, false},
	{"boundary min date", "0001-01-01 00:00", false, false},
	{"boundary max date", "9999-12-31 23:59", false, false},
}

func TestValidateTime(t *testing.T) {
	for _, tt := range validateTimeTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateTime(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("ожидалась ошибка, но err == nil для %q", tt.input)
				}
			}
			if tt.allowZero {
				if !got.IsZero() {
					t.Errorf("ожидалось zeroTime, а получено: %v", got)
				}
			}
		})
	}
}
