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
	tracer := otel.Tracer("db:postgres:categories")
	return &postgresCategoryRepository{conn: conn, tracer: tracer}
}

func (p *postgresCategoryRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Category, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying categories")
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
		cats = append(cats, cat)
	}

	return cats, nil
}

func (p *postgresCategoryRepository) GetByID(ctx context.Context, id int64) (domain.Category, error) {
	query := `
		SELECT
			id,
			name,
			note,
			created_by,
			created_at,
			updated_at
		FROM
			categories
		WHERE
			id = $1
			AND is_deleted = FALSE`

	cat, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.Category{}, err
	}

	if len(cat) == 0 {
		return domain.Category{}, domain.ErrNotFound
	}
	
	return cat[0], nil
}

func (p *postgresCategoryRepository) GetByUserID(ctx context.Context, created_by int64) ([]domain.Category, error) {
	query := `
		SELECT
			id,
			name,
			note,
			created_by,
			created_at,
			updated_at
		FROM
			categories
		WHERE
			created_by = $1
			AND is_deleted = FALSE
		ORDER BY
			name ASC`

	cats, err := p.fetch(ctx, query, created_by)

	if err != nil {
		return []domain.Category{}, err
	}

	return cats, nil
}

// func (p *postgresCategoryRepository) CreateOrUpdate(ctx context.Context, cat *domain.Category) error {
// 	query := `
// 		INSERT INTO categories (name, note, created_by)
// 		VALUES ($1, $2, $3)
// 		ON CONFLICT(username) DO
// 			UPDATE SET access_token = $3,
// 				refresh_token = $4,
// 				token_expires_at = $5,
// 				last_message_id = $6,
// 				is_deleted = FALSE
// 		RETURNING id`

// 	ctx, span := spanWithQuery(ctx, p.tracer, query)
// 	defer span.End()

// 	if err := p.conn.QueryRow(
// 		ctx,
// 		query,
// 		cat.ID,
// 		cat.Name,
// 		cat.Note,
// 	).Scan(&cat.ID); err != nil {
// 		span.SetStatus(codes.Error, "failed upserting category")
// 		span.RecordError(err)
// 		return err
// 	}

// 	return nil
// }

func (p *postgresCategoryRepository) Create(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	query := `
		INSERT INTO categories (name, note, created_by)
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
		span.SetStatus(codes.Error, "failed inserting categories")
		span.RecordError(err)
		return nil, err
	}
	return cat, nil
}

func (p *postgresCategoryRepository) Update(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	query := `
		UPDATE categories
		SET 
			name = $2,
			note = $3,
			updated_at = NOW()
		WHERE 
			id = $1
		RETURNING updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	row := p.conn.QueryRow(
		ctx,
		query,
		cat.ID,
		cat.Name,
		cat.Note,
	); 
	
	if err := row.Scan(&cat.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed to update category")
		span.RecordError(err)
		return nil, err
	}

	return cat, nil
}

func (p *postgresCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE categories
		SET 
			is_deleted = TRUE,
			updated_at = NOW()
		WHERE 
			id = $1`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	result , err := p.conn.Exec(ctx, query, id);

	if err != nil {
		span.SetStatus(codes.Error, "failed to delete category")
		span.RecordError(err)
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}