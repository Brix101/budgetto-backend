package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/Brix101/budgetto-backend/internal/domain"
)

type postgresAccountRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewPostgresAccount(conn Connection) domain.AccountRepository {
	tracer := otel.Tracer("db:postgres:accounts")
	return &postgresAccountRepository{conn: conn, tracer: tracer}
}

func (p *postgresAccountRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Account, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying accounts")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var accs []domain.Account
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(
			&acc.ID,
			&acc.Name,
			&acc.Balance,
			&acc.Note,
			&acc.CreatedBy,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}
	return accs, nil
}

func (p *postgresAccountRepository) GetByID(ctx context.Context, id int64) (domain.Account, error) {
	query := `
		SELECT
			id,
			name,
            balance,
			note,
			created_by,
			created_at,
			updated_at
		FROM 
            accounts
		WHERE 
            id = $1 AND 
            is_deleted = FALSE`

	accs, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Account{}, err
	}

	if len(accs) == 0 {
		return domain.Account{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *postgresAccountRepository) GetByUserID(ctx context.Context, created_by int64) ([]domain.Account, error) {
	query := `
		SELECT
			id,
			name,
            balance,
			note,
			created_by,
			created_at,
			updated_at
		FROM
			accounts
		WHERE
			created_by = $1 AND
			is_deleted = FALSE
		ORDER BY
			name ASC`

	cats, err := p.fetch(ctx, query, created_by)
	if err != nil {
		return []domain.Account{}, err
	}

	if len(cats) <= 0 {
		return []domain.Account{}, nil
	}

	return cats, nil
}

func (p *postgresAccountRepository) Create(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts
			(name, balance, note, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	if err := p.conn.QueryRow(
		ctx,
		query,
		acc.Name,
		acc.Balance,
		acc.Note,
		acc.CreatedBy,
	).Scan(&acc.ID); err != nil {
		span.SetStatus(codes.Error, "failed inserting account")
		span.RecordError(err)
		return err
	}

	return nil
}

func (p *postgresAccountRepository) Update(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	query := `
		UPDATE accounts
		SET 
			name = $2,
			balance = $3,
			note = $4,
			updated_at = NOW()
		WHERE 
			id = $1
		RETURNING updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	row := p.conn.QueryRow(
		ctx,
		query,
		acc.ID,
		acc.Name,
		acc.Balance,
		acc.Note,
	)

	if err := row.Scan(&acc.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed to update account")
		span.RecordError(err)
		return nil, err
	}

	return acc, nil
}

func (p *postgresAccountRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE accounts
		SET 
			is_deleted = TRUE,
			updated_at = NOW()
		WHERE 
			id = $1`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	result, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		span.SetStatus(codes.Error, "failed to delete account")
		span.RecordError(err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
