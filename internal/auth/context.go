package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const userIdKey = "userId"

func (a *Auth) GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userId := ctx.Value(userIdKey)
	if userId == nil {
		return uuid.Nil, errors.New("userId not found in context")
	}
	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("userId not in correct format")
	}
	return userUUID, nil
}

func (a *Auth) PutUserIdToContext(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}
