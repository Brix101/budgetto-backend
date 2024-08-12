package domain

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
	Sub   int    `json:"sub"`
}

type userToken struct {
	AccessToken string `json:"access_token"`
}

func (u User) GenerateClaims() (string, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET")

	claims := UserClaims{
		jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", u.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		u.Name,
		u.Email,
		int(u.ID),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
