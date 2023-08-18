package repository

import (
	"context"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgresCategoryRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewPostgresCategory(conn Connection) domain.CategoryRepository {
	tracer := otel.Tracer("db:postgres:accounts")
	return &postgresCategoryRepository{conn: conn, tracer: tracer}
}

func (p *postgresCategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Category, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying accounts")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var accs []domain.Category
	for rows.Next() {
		var acc domain.Category
		if err := rows.Scan(
			&acc.ID,
			&acc.Name,
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

func (p *postgresCategoryRepository) GetByID(ctx context.Context, id int64) (domain.Category, error) {
	query := `
		SELECT *
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL`

	accs, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Category{}, err
	}

	if len(accs) == 0 {
		return domain.Category{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *postgresCategoryRepository) GetByUserID(ctx context.Context, id int64) ([]domain.Category, error) {
	query := `
		SELECT *
		FROM categories
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC`

	return p.fetch(ctx, query, id)
}
