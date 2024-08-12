package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Brix101/budgetto-backend/internal/domain"
)

const BudgettoToken = "budgetto-token"

type UserCtxKey struct{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenSecret := os.Getenv("TOKEN_SECRET")
		cookieToken, err := r.Cookie(BudgettoToken)
		if err == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookieToken.Value, func(token *jwt.Token) (interface{}, error) {
			// You should implement your own logic to validate the token and return the appropriate key
			// For example, you could use a secret key or a public key
			return []byte(tokenSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set the claims in the context
		claims, err := transformMapClaimsToUserClaims(token.Claims)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)

	return token
}

func transformMapClaimsToUserClaims(claims jwt.Claims) (*domain.UserClaims, error) {
	fmt.Println(claims.GetSubject())

	if jwtClaims, ok := claims.(jwt.MapClaims); ok {
		sub, ok := jwtClaims["sub"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid 'sub' claim")
		}

		name, ok := jwtClaims["name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'name' claim")
		}

		email, ok := jwtClaims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid 'email' claim")
		}

		// Create a new instance of *UserClaims with the extracted values
		userClaims := &domain.UserClaims{
			Sub:   int(sub),
			Name:  name,
			Email: email,
		}

		// Add other custom claim fields here if needed
		return userClaims, nil
	}

	return nil, fmt.Errorf("failed to transform jwt.MapClaims to *UserClaims")
}
