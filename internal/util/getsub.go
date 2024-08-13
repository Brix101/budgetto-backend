package util

import (
	"context"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/golang-jwt/jwt/v5"
)

func GetSub(ctx context.Context) (uint, error) {
	user := ctx.Value(middlewares.UserCtxKey{}).(jwt.MapClaims)
	sub, err := user.GetSubject()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}
