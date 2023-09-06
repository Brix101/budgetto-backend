package repository

import (
	"context"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgresTransactionRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewPostgresTransaction(conn Connection) domain.TransactionRepository {
	tracer := otel.Tracer("db:postgres:transactions")
	return &postgresTransactionRepository{conn: conn, tracer: tracer}
}

func (p *postgresTransactionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Transaction, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying transactions")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	trns := []domain.Transaction{}
	for rows.Next() {
		var trn domain.Transaction
		var acc domain.Account
		var cat domain.Category
		if err := rows.Scan(
			&trn.ID,
			&trn.Amount,
			&trn.Note,
			&trn.Operation,
			&trn.AccountID,
			&trn.CategoryID,
			&trn.CreatedBy,
			&trn.CreatedAt,
			&trn.UpdatedAt,
			&acc.ID,
			&acc.Name,
			&acc.Balance,
			&acc.Note,
			&acc.CreatedAt,
			&acc.UpdatedAt,
			&cat.ID,
			&cat.Name,
			&cat.Note,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		); err != nil {
			return nil, err
		}

		trn.Account = acc
		trn.Category = cat
		trns = append(trns, trn)
	}
	return trns, nil
}

func (p *postgresTransactionRepository) GetByID(ctx context.Context, id int64) (domain.Transaction, error) {
	query := `
		SELECT 
			T.ID,
			T.amount,
			T.note,
			T.operation,
			T.account_id,
			T.category_id,
			T.created_by,
			T.created_at,
			T.updated_at,
			A.ID AS acc_id,
			A.NAME AS acc_name,
			A.balance AS acc_balance,
			A.note AS acc_note,
			A.created_at AS acc_created_at,
			A.updated_at AS acc_updated_at,
			C.ID AS cat_id,
			C.NAME AS cat_name,
			C.note AS cat_note,
			C.created_at AS cat_created_at,
			C.updated_at AS cat_updated_at 
		FROM
			transactions T 
			JOIN accounts A ON T.account_id = A.ID 
			JOIN categories C ON T.category_id = C.ID 
		WHERE
			T.ID = $1 
			AND T.is_deleted = FALSE;`

	trns, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(trns) == 0 {
		return domain.Transaction{}, domain.ErrNotFound
	}
	return trns[0], nil
}

func (p *postgresTransactionRepository) GetByUserSUB(ctx context.Context, sub string) ([]domain.Transaction, error) {
	query := `
		SELECT 
			T.ID,
			T.amount,
			T.note,
			T.operation,
			T.account_id,
			T.category_id,
			T.created_by,
			T.created_at,
			T.updated_at,
			A.ID AS acc_id,
			A.NAME AS acc_name,
			A.balance AS acc_balance,
			A.note AS acc_note,
			A.created_at AS acc_created_at,
			A.updated_at AS acc_updated_at,
			C.ID AS cat_id,
			C.NAME AS cat_name,
			C.note AS cat_note,
			C.created_at AS cat_created_at,
			C.updated_at AS cat_updated_at 
		FROM
			transactions T 
			JOIN accounts A ON T.account_id = A.ID 
			JOIN categories C ON T.category_id = C.ID 
		WHERE
			T.created_by = $1 
			AND T.is_deleted = FALSE;`

	trns, err := p.fetch(ctx, query, sub)
	if err != nil {
		return []domain.Transaction{}, err
	}

	return trns, nil
}

func (p *postgresTransactionRepository) GetOperationType(ctx context.Context) ([]string, error) {
	query := `
        SELECT enumlabel
        FROM pg_enum
        WHERE enumtypid = 'operation'::regtype;`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying operations")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	enumLabels := []string{}
	for rows.Next() {
		var enumLabel string
		if err := rows.Scan(
			&enumLabel,
		); err != nil {
			return nil, err
		}
		enumLabels = append(enumLabels, enumLabel)
	}

	return enumLabels, nil
}

func (p *postgresTransactionRepository) Create(ctx context.Context, trn *domain.Transaction) (*domain.Transaction, error) {
	query := `
		INSERT INTO transactions
			(amount, note, account_id, category_id, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	if err := p.conn.QueryRow(
		ctx,
		query,
		trn.Amount,
		trn.CategoryID,
		trn.CreatedBy,
	).Scan(
		&trn.ID,
		&trn.CreatedAt,
		&trn.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed inserting transaction")
		span.RecordError(err)
		return nil, err
	}

	return trn, nil
}

func (p *postgresTransactionRepository) Update(ctx context.Context, trn *domain.Transaction) (*domain.Transaction, error) {
	query := `
		UPDATE transactions
		SET 
			amount = $2,
			note = $3,
			account_id = $4,
			category_id = $5,
			updated_at = NOW()
		WHERE 
			id = $1
		RETURNING updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	row := p.conn.QueryRow(
		ctx,
		query,
		trn.ID,
		trn.Amount,
		trn.Note,
		trn.AccountID,
		trn.CategoryID,
	)

	if err := row.Scan(&trn.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed to update transaction")
		span.RecordError(err)
		return nil, err
	}

	return trn, nil
}

func (p *postgresTransactionRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE transactions
		SET 
			is_deleted = TRUE,
			updated_at = NOW()
		WHERE 
			id = $1`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	result, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		span.SetStatus(codes.Error, "failed to delete transaction")
		span.RecordError(err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
