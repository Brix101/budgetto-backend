package middlwares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Brix101/budgetto-backend/config"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func Auth0Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := config.GetConfig()
		issuerURL, err := url.Parse("https://" + env.AUTH0_DOMAIN + "/")
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}
		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{env.AUTH0_AUDIENCE},
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}
		fmt.Println("111111111111111111111111111111111111111111111111111111111111111111111")
		// get the token from the request header
		authHeader := r.Header.Get("Authorization")
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println("222222222222222222222222222222222222222222222222222222222222222222222")
		// Validate the token
		// tokenInfo, err := jwtValidator.ValidateToken(r.Context(), authHeaderParts[1])
		// if err != nil {
		// 	fmt.Println(err)
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		fmt.Println(authHeaderParts[1], jwtValidator)
		next.ServeHTTP(w, r)
	})
}
