package domain

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var AccessExp = time.Hour * 1
var RefreshExp = time.Hour * 24 * 90

type UserClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) GenerateClaims() (string, error) {
	privateKey := os.Getenv("ACCESS_PRIVATE_KEY")
	keyData, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData))
	if err != nil {
		return "", err
	}

	claims := UserClaims{
		jwt.RegisteredClaims{
			Issuer:    "Budgetto",
			ID:        fmt.Sprintf("%d", u.ID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessExp)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", u.ID),
		},
		u.Name,
		u.Email,
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(parsedKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (u User) GenerateRefreshToken() (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "Budgetto",
		ID:        fmt.Sprintf("%d", u.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshExp)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   fmt.Sprintf("%d", u.ID),
	}

	privateKey := os.Getenv("REFRESH_PRIVATE_KEY")
	keyData, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyData))
	if err != nil {
		return "", err
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(parsedKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

type UserWithToken struct {
	User        User   `json:"user"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

func (u User) GenerateUserWithToken() (*UserWithToken, error) {
	accessToken, err := u.GenerateClaims()
	if err != nil {
		return nil, err
	}

	expirationTime := time.Now().Add(AccessExp)
	expiresIn := int64(time.Until(expirationTime).Seconds()) + 1

	return &UserWithToken{
		User:        u,
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}, nil
}
