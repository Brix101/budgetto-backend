package api

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Brix101/budgetto-backend/internal/domain"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"go.uber.org/zap"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}

	return false
}

func (a api) authClaims(ctx context.Context) (*CustomClaims, error) {
	token, err := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !err {
		return nil, domain.ErrForbidden
	}

	claims := token.CustomClaims.(*CustomClaims)
	if claims.Sub == "" {
		return nil, domain.ErrForbidden
	}

	return claims, nil
}

func (a api) auth0Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domain := os.Getenv("AUTH0_DOMAIN")
		authdience := os.Getenv("AUTH0_AUDIENCE")

		issuerURL, err := url.Parse("https://" + domain + "/")
		if err != nil {
			a.logger.Error("Failed to parse the issuer url", zap.Error(err))
		}
		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{authdience},
			validator.WithCustomClaims(
				func() validator.CustomClaims {
					return &CustomClaims{}
				},
			),
			validator.WithAllowedClockSkew(time.Minute),
		)
		if err != nil {
			a.logger.Error("Failed to set up the jwt validator", zap.Error(err))
		}

		// get the token from the request header
		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, _ error) {
			// a.logger.Info("Invalid JWT", zap.Error(err))

			// w.Header().Set("Content-Type", "application/json")
			// w.Header().Set("X-Budgetto-Error", err.Error())
			// w.WriteHeader(http.StatusUnauthorized)
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)

			next.ServeHTTP(w, r)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}
