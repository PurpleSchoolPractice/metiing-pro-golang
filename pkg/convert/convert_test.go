package convert_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/convert"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestParseId(t *testing.T) {
	// Максимальное значение uint для текущей архитектуры.
	maxUint := ^uint(0)
	maxUintString := strconv.FormatUint(uint64(maxUint), 10)

	// Значение на единицу больше максимального uint.
	overflowString := "4294967296"
	if strconv.IntSize == 64 {
		overflowString = "18446744073709551616"
	}

	tests := []struct {
		name        string
		rawID       string
		addParam    bool
		expectedID  uint
		expectError bool
	}{
		{
			name:        "valid ID",
			rawID:       "42",
			addParam:    true,
			expectedID:  42,
			expectError: false,
		},
		{
			name:        "minimum valid ID",
			rawID:       "1",
			addParam:    true,
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "maximum uint",
			rawID:       maxUintString,
			addParam:    true,
			expectedID:  maxUint,
			expectError: false,
		},
		{
			name:        "leading zeroes",
			rawID:       "001",
			addParam:    true,
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "missing parameter",
			addParam:    false,
			expectError: true,
		},
		{
			name:        "empty parameter",
			rawID:       "",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "zero",
			rawID:       "0",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "negative number",
			rawID:       "-1",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "plus sign",
			rawID:       "+1",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "not a number",
			rawID:       "abc",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "fractional number",
			rawID:       "1.5",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "number with spaces",
			rawID:       " 1 ",
			addParam:    true,
			expectError: true,
		},
		{
			name:        "uint overflow",
			rawID:       overflowString,
			addParam:    true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// URL-параметры Chi хранятся в RouteContext запроса.
			routeContext := chi.NewRouteContext()
			if tt.addParam {
				routeContext.URLParams.Add("id", tt.rawID)
			}

			ctx := context.WithValue(
				req.Context(),
				chi.RouteCtxKey,
				routeContext,
			)
			req = req.WithContext(ctx)

			actualID, err := convert.ParseId(req, "id")

			if tt.expectError {
				require.ErrorIs(t, err, convert.ErrInvalidID)
				require.Zero(t, actualID)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedID, actualID)
		})
	}
}
