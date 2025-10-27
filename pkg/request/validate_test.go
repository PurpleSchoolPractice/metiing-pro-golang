package request_test

import (
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/request"
	"github.com/stretchr/testify/require"
)

type validateTimeCase struct {
	name           string
	input          string
	expectError    bool
	expectZeroTime bool
}

var testCases = []validateTimeCase{
	{
		name:           "valid date",
		input:          "2025-10-31 14:30",
		expectError:    false,
		expectZeroTime: false,
	},
	{
		name:           "invalid format",
		input:          "25-10-2025 14:30",
		expectError:    true,
		expectZeroTime: true,
	},
	{
		name:           "invalid date values",
		input:          "2025-13-32 14:30",
		expectError:    true,
		expectZeroTime: true,
	},
	{
		name:           "boundary min date",
		input:          "0001-01-01 00:00",
		expectError:    false,
		expectZeroTime: true,
	},
	{
		name:           "boundary max date",
		input:          "9999-12-31 23:59",
		expectError:    false,
		expectZeroTime: false,
	},
}

func TestValidateTime(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := request.ValidateTime(tc.input)

			if tc.expectError {
				require.Error(t, err, "ожидалась ошибка, но err == nil для %q", tc.input)
			} else {
				require.NoError(t, err, "не ожидалась ошибка, но err != nil для %q: %v", tc.input, err)
			}

			if tc.expectZeroTime {
				require.True(t, got.IsZero(), "ожидалось zeroTime, а получено: %v", got)
			} else {
				require.False(t, got.IsZero(), "не ожидалось zeroTime, но получили zeroTime для %q", tc.input)
			}
		})
	}
}
