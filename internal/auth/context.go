package auth

import (
	"context"
	"errors"
)

const userIdKey = "userId"

func (a *Auth) GetUserIdFromContext(ctx context.Context) (string, error) {
	userId := ctx.Value(userIdKey)
	if userId == nil {
		return "", errors.New("userId not found in context")
	}
	return userId.(string), nil
}

func (a *Auth) PutUserIdToContext(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}
