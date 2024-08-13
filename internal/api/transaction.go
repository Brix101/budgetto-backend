package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"github.com/Brix101/budgetto-backend/internal/middlewares"
	"github.com/Brix101/budgetto-backend/internal/util"
)

type TransactionCtx struct{}

func (a api) TransactionRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.Auth)

	r.Get("/", a.transactionListHandler)
	r.Post("/", a.transactionCreateHandler)
	r.Get("/operations", a.transactionOpListHandler)

	r.Route("/{id}", func(r chi.Router) {
		r.Use(a.TransctionCtx)

		r.Get("/", a.transactionGetHandler)
		r.Put("/", a.transactionUpdateHandler)
		r.Delete("/", a.transactionDeleteHandler)
	})

	return r
}

func (a api) TransctionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sub, err := util.GetSub(ctx)
		if err != nil {
			a.errorResponse(w, r, 500, err)
			return
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			a.errorResponse(w, r, 500, err)
			return
		}

		item, err := a.transactionRepo.GetByID(ctx, uint(id))
		if err != nil {
			status := 500
			if err.Error() == domain.ErrNotFound.Error() {
				status = 404
			}
			a.errorResponse(w, r, status, err)
			return
		}

		if item.CreatedBy != sub {
			a.errorResponse(w, r, 403, domain.ErrForbidden)
			return
		}

		ctx = context.WithValue(ctx, CatCtx{}, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type createTransactionRequest struct {
	Amount     float64 `json:"amount" validate:"gte=0"`
	Note       string  `json:"note"`
	Operation  string  `json:"operation" validate:"oneof=Expense Income Transfer Refund"`
	AccountID  uint    `json:"account_id" validate:"required"`
	CategoryID uint    `json:"category_id" validate:"required"`
}

func (a api) transactionOpListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	ops, err := a.transactionRepo.GetOperationType(ctx)
	if err != nil {
		a.logger.Error("failed to fetch operations from database", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(ops)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) transactionListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := ctx.Value(middlewares.UserCtxKey{}).(*domain.UserClaims)
	sub, err := user.GetSubject()
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	trns, err := a.transactionRepo.GetByUserSUB(ctx, sub)
	if err != nil {
		a.logger.Error("failed to fetch transactions from database", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(trns)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) transactionCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	sub, err := util.GetSub(ctx)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	reqBody := createTransactionRequest{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		a.logger.Error("failed to validate create budget struct", zap.Error(err))
		a.errorResponse(w, r, 400, err)
		return
	}

	trnReq := domain.Transaction{
		Amount:     reqBody.Amount,
		CategoryID: reqBody.CategoryID,
		AccountID:  reqBody.AccountID,
		CreatedBy:  sub,
	}

	newTrn, err := a.transactionRepo.Create(ctx, &trnReq)
	if err != nil {
		a.logger.Error("failed to create transaction", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, newTrn.ID)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(trn)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) transactionGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(TransactionCtx{}).(domain.Transaction)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	resJSON, err := json.Marshal(item)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) transactionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(TransactionCtx{}).(domain.Transaction)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	upTrn, err := a.transactionRepo.Update(ctx, &item)
	if err != nil {
		a.logger.Error("failed to update transaction", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, upTrn.ID)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	resJSON, err := json.Marshal(trn)
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (a api) transactionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	item, ok := ctx.Value(TransactionCtx{}).(domain.Transaction)
	if !ok {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	if err := a.transactionRepo.Delete(ctx, item.ID); err != nil {
		a.logger.Error("failed to delete transactiond", zap.Error(err))
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
	}

	data := map[string]string{
		"message": "Item deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
