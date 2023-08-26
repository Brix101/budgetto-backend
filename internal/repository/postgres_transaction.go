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
		if err := rows.Scan(
			&trn.ID,
			&trn.Amount,
			&trn.Note,
			&trn.TransactionType,
			&trn.AccountID,
			&trn.CategoryID,
			&trn.CreatedBy,
			&trn.CreatedAt,
			&trn.UpdatedAt,
		); err != nil {
			return nil, err
		}
		trns = append(trns, trn)
	}
	return trns, nil
}

func (p *postgresTransactionRepository) GetByID(ctx context.Context, id int64) (domain.Transaction, error) {
	query := `
		SELECT
			id,
            amount,
            note,
            transaction_type,
            account_id,
			category_id,
			created_by,
			created_at,
			updated_at
		FROM 
            transactions
		WHERE 
            id = $1 AND 
            is_deleted = FALSE`

	trns, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(trns) == 0 {
		return domain.Transaction{}, domain.ErrNotFound
	}
	return trns[0], nil
}

func (p *postgresTransactionRepository) GetByUserID(ctx context.Context, created_by int64) ([]domain.Transaction, error) {
	query := `
		SELECT
			id,
            amount,
            note,
            transaction_type,
            account_id,
			category_id,
			created_by,
			created_at,
			updated_at
		FROM
			transactions
		WHERE
			created_by = $1 AND
			is_deleted = FALSE
		ORDER BY
			name ASC`

	trns, err := p.fetch(ctx, query, created_by)
	if err != nil {
		return []domain.Transaction{}, err
	}

	return trns, nil
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
