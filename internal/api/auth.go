package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/Brix101/budgetto-backend/internal/util"
)

func (a api) AuthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/refresh-token", a.refreshHandler)
	r.Post("/sign-in", a.signInHandler)
	r.Post("/sign-up", a.signUpHandler)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.Auth)
		r.Get("/me", a.meHandler)
	})

	return r
}

type signInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"` // Minimum length: 6
}

type signUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"` // Minimum length: 6
}

func (a api) signInHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var reqBody signInRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	usr, err := a.userRepo.GetByEmail(ctx, reqBody.Email)
	if err != nil {
		a.errorResponse(w, r, 401, domain.ErrInvalidCredentials)
		return
	}

	if validatePass := usr.CheckPassword(reqBody.Password); !validatePass {
		a.errorResponse(w, r, 401, domain.ErrInvalidCredentials)
		return
	}

	token, err := usr.GenerateRefreshToken()
	if err != nil {
		a.logger.Error("failed to generate user claims", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	data, err := usr.GenerateUserWithToken()
	if err != nil {
		a.logger.Error("failed to generate user with token", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(data)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	// Create and set cookies in the response
	cookie := http.Cookie{
		Name:     middlewares.BudgetttoCookieKey, // Cookie name
		Value:    token,                          // Cookie value (you can customize this)
		Path:     "/",                            // Cookie path
		HttpOnly: true,                           // Prevent JavaScript access
		Expires:  time.Now().Add(domain.RefreshExp),
		// You can set more attributes like Expires, MaxAge, Secure, etc. as needed.
	}

	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) signUpHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	reqBody := signUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	newUsr := domain.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	}

	if err := newUsr.HashPassword(); err != nil {
		a.errorResponse(w, r, 400, err)
		return
	}

	usr, err := a.userRepo.Create(ctx, &newUsr)
	if err != nil {
		a.logger.Error("failed to create user", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(usr)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) refreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	cookie, err := r.Cookie(middlewares.BudgetttoCookieKey)
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	publicKey := os.Getenv("REFRESH_PUBLIC_KEY")
	keyData, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyData))
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return parsedKey, nil
	})

	if err != nil || !token.Valid {
		a.errorResponse(w, r, 401, err)
		return
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	usr, err := a.userRepo.GetByID(ctx, uint(userId))
	if err != nil {
		a.errorResponse(w, r, 401, err)
		return
	}

	data, err := usr.GenerateUserWithToken()
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(data)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) meHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	sub, err := util.GetSub(ctx)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	user, err := a.userRepo.GetByID(ctx, sub)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(user)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}
