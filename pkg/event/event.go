package event

import "context"

// GetUserIDFromContext извлекает ID пользователя из контекста запроса
func GetUserIDFromContext(ctx context.Context) uint {
	userID, ok := ctx.Value("userID").(uint)
	if !ok {
		return 0 // Возвращаем 0, если пользователь не найден в контексте
	}
	return userID
}
