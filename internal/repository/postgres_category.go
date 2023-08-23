package repository

import (
	"context"
	"log"

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
	tracer := otel.Tracer("db:postgres:catounts")
	return &postgresCategoryRepository{conn: conn, tracer: tracer}
}

func (p *postgresCategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Category, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying catounts")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Note,
			&cat.CreatedBy,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		); err != nil {
			return nil, err
		}

		log.Println(cat)
		cats = append(cats, cat)
	}
	return cats, nil
}

func (p *postgresCategoryRepository) GetByID(ctx context.Context, id int64) (domain.Category, error) {
	log.Println(id)
	query := `
		SELECT *
		FROM categories
		WHERE id = $1 
		AND is_deleted = false`

	cats, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Category{}, err
	}

	if len(cats) == 0 {
		return domain.Category{}, domain.ErrNotFound
	}
	return cats[0], nil
}

func (p *postgresCategoryRepository) GetByUserID(ctx context.Context, id int64) ([]domain.Category, error) {
	query := `
		SELECT *
		FROM categories
		WHERE user_id = $1 
		OR is_deleted = false
		ORDER BY name ASC`

	return p.fetch(ctx, query, id)
}

func (p *postgresCategoryRepository) Create(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	query := `
		INSERT INTO categories (name, note, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	if err := p.conn.QueryRow(
		ctx,
		query,
		cat.Name,
		cat.Note,
		cat.CreatedBy,
	).Scan(
		&cat.ID,
		&cat.CreatedAt,
		&cat.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed inserting catount")
		span.RecordError(err)
		return nil, err
	}
	return cat, nil
}
