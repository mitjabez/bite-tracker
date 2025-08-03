package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitjabez/bite-tracker/internal/model"
)

const userKey = "user"

func (a *Auth) GetUserFromContext(ctx context.Context) (model.User, error) {
	userAny := ctx.Value(userKey)
	if userAny == nil {
		return model.User{}, errors.New("user not found in context")
	}
	user, ok := userAny.(model.User)
	if !ok {
		return model.User{}, fmt.Errorf("user not in correct format: %v", userAny)
	}
	return user, nil
}

func (a *Auth) PutUserToContext(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
