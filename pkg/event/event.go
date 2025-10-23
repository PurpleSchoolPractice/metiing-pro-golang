package event

import (
	"context"
	"errors"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
)

var EventErrors = map[string]error{
	"type":    errors.New("INVALID_TYPE_ERR"),
	"missing": errors.New("MISSING_CTX_ERR"),
}

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	if v := ctx.Value(middleware.ContextUserIDKey); v != nil {
		if id, ok := v.(uint); ok {
			return id, nil
		} else {
			return 0, EventErrors["type"]
		}
	}

	return 0, EventErrors["missing"]
}
