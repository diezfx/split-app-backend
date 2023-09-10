package contextutil

import (
	"context"
)

type ctxKey string

const userIDKey ctxKey = "userId"

func GetUserIDFromCtx(ctx context.Context) string {
	maybeUser := ctx.Value(userIDKey)

	user, ok := maybeUser.(string)
	if !ok {
		return ""
	}
	return user
}

func AddUserIDToCtx(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, userIDKey, value)
}
