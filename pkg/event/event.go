package event

import (
	"context"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
)

func GetUserIDFromContext(ctx context.Context) uint {
	if v := ctx.Value(middleware.ContextUserIDKey); v != nil {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}
