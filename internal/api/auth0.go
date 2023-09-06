package api

import (
	"context"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)


func (a api) authClaims(ctx context.Context)(*middlewares.CustomClaims, error) {
	token, err := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !err {
		return nil, domain.ErrForbidden
	}

	claims := token.CustomClaims.(*middlewares.CustomClaims)
	if claims.Sub == "" {
		return nil, domain.ErrForbidden
	}

	return claims, nil
}
