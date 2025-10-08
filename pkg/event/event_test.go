package event_test

import (
	"context"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	context     context.Context
	expected    any
	expectError bool
}{
	{context: context.WithValue(context.Background(), middleware.ContextUserIDKey, uint(42)), expected: uint(42), expectError: false},
	{context: context.WithValue(context.Background(), middleware.ContextUserIDKey, int(42)), expected: uint(0), expectError: true},
	{context: context.WithValue(context.Background(), middleware.ContextUserIDKey, nil), expected: uint(0), expectError: true},
	{context: context.WithValue(context.Background(), middleware.ContextUserIDKey, "42"), expected: uint(0), expectError: true},
	{context: context.WithValue(context.Background(), middleware.ContextUserIDKey, struct{}{}), expected: uint(0), expectError: true},
}

func TestGetUserIDFromContext(t *testing.T) {
	for _, testData := range testCases {
		t.Run("GetUserIDFromContext", func(t *testing.T) {
			res, err := event.GetUserIDFromContext(testData.context)
			if testData.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, testData.expected, res)
		})
	}

}
