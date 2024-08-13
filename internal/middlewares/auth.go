package middlewares

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const BudgetttoCookieKey = "x-budgetto-token"

type AuthCtx struct{}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authToken := strings.Replace(authHeader, "Bearer ", "", 1)

		publicKey := os.Getenv("ACCESS_PUBLIC_KEY")
		keyData, err := base64.StdEncoding.DecodeString(publicKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyData))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			// You should implement your own logic to validate the token and return the appropriate key
			// For example, you could use a secret key or a public key
			return parsedKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthCtx{}, token.Claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
