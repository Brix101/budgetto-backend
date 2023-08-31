package repository

import (
	"context"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgresBudgetRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewPostgresBudget(conn Connection) domain.BudgetRepository {
	tracer := otel.Tracer("db:postgres:budgets")
	return &postgresBudgetRepository{conn: conn, tracer: tracer}
}

func (p *postgresBudgetRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Budget, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying budgets")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	buds := []domain.Budget{}
	for rows.Next() {
		var bud domain.Budget
		var cat domain.Category
		if err := rows.Scan(
			&bud.ID,
			&bud.Amount,
			&bud.CategoryID,
			&bud.CreatedBy,
			&bud.CreatedAt,
			&bud.UpdatedAt,
			&cat.ID,
			&cat.Name,
			&cat.Note,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		); err != nil {
			return nil, err
		}
		bud.Category = cat
		buds = append(buds, bud)
	}
	return buds, nil
}

func (p *postgresBudgetRepository) GetByID(ctx context.Context, id int64) (domain.Budget, error) {
	query := `
		SELECT
	        b.id,
	        b.amount,
	        b.category_id,
	        b.created_by,
	        b.created_at,
	        b.updated_at,
	        c.id AS category_id,
	        c.name AS category_name,
	        c.note AS category_note,    
	        c.created_at AS cat_created_at, 
	        c.updated_at AS cat_updated_at
        FROM
	        budgets b
	    JOIN categories c ON b.category_id = c.ID 
        WHERE
	        b.id = $1 
	        AND b.is_deleted = FALSE;`

	buds, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Budget{}, err
	}

	if len(buds) == 0 {
		return domain.Budget{}, domain.ErrNotFound
	}
	return buds[0], nil
}

func (p *postgresBudgetRepository) GetByUserID(ctx context.Context, created_by int64) ([]domain.Budget, error) {
	query := `
        SELECT
	        b.id,
	        b.amount,
	        b.category_id,
	        b.created_by,
	        b.created_at,
	        b.updated_at,
	        c.id AS category_id,
	        c.name AS category_name,
	        c.note AS category_note,    
	        c.created_at AS cat_created_at, 
	        c.updated_at AS cat_updated_at
        FROM
	        budgets b
	    JOIN categories c ON b.category_id = c.ID 
        WHERE
	        b.created_by = $1 
	        AND b.is_deleted = FALSE;
	    ORDER BY
		    c.name ASC`

	buds, err := p.fetch(ctx, query, created_by)
	if err != nil {
		return []domain.Budget{}, err
	}

	return buds, nil
}

func (p *postgresBudgetRepository) Create(ctx context.Context, bud *domain.Budget) (*domain.Budget, error) {
	query := `
		INSERT INTO budgets
			(amount, category_id, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	if err := p.conn.QueryRow(
		ctx,
		query,
		bud.Amount,
		bud.CategoryID,
		bud.CreatedBy,
	).Scan(
		&bud.ID,
		&bud.CreatedAt,
		&bud.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed inserting budget")
		span.RecordError(err)
		return nil, err
	}

	return bud, nil
}

func (p *postgresBudgetRepository) Update(ctx context.Context, bud *domain.Budget) (*domain.Budget, error) {
	query := `
		UPDATE budgets
		SET 
			amount = $2,
			category_id = $3,
			updated_at = NOW()
		WHERE 
			id = $1
		RETURNING updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	row := p.conn.QueryRow(
		ctx,
		query,
		bud.ID,
		bud.Amount,
		bud.CategoryID,
	)

	if err := row.Scan(&bud.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed to update budget")
		span.RecordError(err)
		return nil, err
	}

	return bud, nil
}

func (p *postgresBudgetRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE budgets
		SET 
			is_deleted = TRUE,
			updated_at = NOW()
		WHERE 
			id = $1`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	result, err := p.conn.Exec(ctx, query, id)
	if err != nil {
		span.SetStatus(codes.Error, "failed to delete budget")
		span.RecordError(err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
