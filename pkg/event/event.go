package event

import (
	"context"
	"errors"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
)

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	if v := ctx.Value(middleware.ContextUserIDKey); v != nil {
		if id, ok := v.(uint); ok {
			return id, nil
		}
	}

	return 0, errors.New("user ID in context has invalid type")
}
