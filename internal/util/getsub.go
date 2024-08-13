package util

import (
	"context"
	"strconv"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
)

func GetSub(ctx context.Context) (uint, error) {
	user := ctx.Value(middlewares.UserCtxKey{}).(*domain.UserClaims)
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
