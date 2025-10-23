package event_test

import (
	"context"
	"testing"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/stretchr/testify/require"
)

var testCases = []struct {
	name        string
	context     context.Context
	expected    any
	expectError error
}{
	{
		name:        "Valid uint user ID",
		context:     context.WithValue(context.Background(), middleware.ContextUserIDKey, uint(42)),
		expected:    uint(42),
		expectError: nil,
	},
	{
		name:        "User ID wrong type int",
		context:     context.WithValue(context.Background(), middleware.ContextUserIDKey, int(42)),
		expected:    uint(0),
		expectError: event.EventErrors["type"],
	},
	{
		name:        "User ID missing",
		context:     context.WithValue(context.Background(), middleware.ContextUserIDKey, nil),
		expected:    uint(0),
		expectError: event.EventErrors["missing"],
	},
	{
		name:        "User ID wrong type string",
		context:     context.WithValue(context.Background(), middleware.ContextUserIDKey, "42"),
		expected:    uint(0),
		expectError: event.EventErrors["type"],
	},
	{
		name:        "User ID wrong type struct",
		context:     context.WithValue(context.Background(), middleware.ContextUserIDKey, struct{}{}),
		expected:    uint(0),
		expectError: event.EventErrors["type"],
	},
}

func TestGetUserIDFromContext(t *testing.T) {
	for _, testData := range testCases {
		t.Run(testData.name, func(t *testing.T) {
			res, err := event.GetUserIDFromContext(testData.context)
			require.ErrorIs(t, err, testData.expectError)
			require.Equal(t, testData.expected, res)
		})
	}
}
