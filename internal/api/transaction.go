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
)

func (a api) TransactionRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.AuthMiddleware)

	r.Get("/", a.transactionListHandler)
	r.Post("/", a.transactionCreateHandler)
	r.Get("/{id}", a.transactionGetHandler)
	r.Put("/{id}", a.transactionUpdateHandler)
	r.Delete("/{id}", a.transactionDeleteHandler)
	r.Get("/operations", a.transactionOpListHandler)

	return r
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

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

	trns, err := a.transactionRepo.GetByUserSUB(ctx, int64(user.Sub))
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

func (a api) transactionGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		if status == 500 {
			a.logger.Error("failed to fetch from database", zap.Error(err))
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if trn.CreatedBy != user.Sub {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
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

func (a api) transactionCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

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
		CreatedBy:  user.Sub,
	}

	newTrn, err := a.transactionRepo.Create(ctx, &trnReq)
	if err != nil {
		a.logger.Error("failed to create transaction", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, int64(newTrn.ID))
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

func (a api) transactionUpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	reqTrn, err := a.transactionRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if reqTrn.CreatedBy != user.Sub {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&reqTrn); err != nil {
		a.logger.Error("failed to parse request json", zap.Error(err))
		a.errorResponse(w, r, 422, err)
		return
	}
	defer r.Body.Close()

	upTrn, err := a.transactionRepo.Update(ctx, &reqTrn)
	if err != nil {
		a.logger.Error("failed to update transaction", zap.Error(err))
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, int64(upTrn.ID))
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

	user := r.Context().Value(middlewares.UserCtxKey{}).(*domain.UserClaims)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		a.errorResponse(w, r, 500, err)
		return
	}

	trn, err := a.transactionRepo.GetByID(ctx, int64(id))
	if err != nil {
		status := 500
		if err.Error() == domain.ErrNotFound.Error() {
			status = 404
		}
		a.errorResponse(w, r, status, err)
		return
	}

	if trn.CreatedBy != user.Sub {
		a.errorResponse(w, r, 403, domain.ErrForbidden)
		return
	}

	if err := a.transactionRepo.Delete(ctx, int64(id)); err != nil {
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
