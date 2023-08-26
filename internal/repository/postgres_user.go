package repository

import (
	"context"

	"github.com/Brix101/budgetto-backend/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type postgresUserRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewPostgresUser(conn Connection) domain.UserRepository {
	tracer := otel.Tracer("db:postgres:users")
	return &postgresUserRepository{conn: conn, tracer: tracer}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying users")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	usrs := []domain.User{}
	for rows.Next() {

		var usr domain.User
		if err := rows.Scan(
			&usr.ID,
			&usr.Name,
			&usr.Email,
			&usr.Password,
			&usr.Bio,
			&usr.Image,
			&usr.CreatedAt,
			&usr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		usrs = append(usrs, usr)
	}

	return usrs, nil
}

func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	query := `
		SELECT
			id,
			name,
			email,
			password,
			bio,
			image,
			created_at,
			updated_at
		FROM
			users
		WHERE
			id = $1
			AND is_deleted = FALSE`

	usr, err := p.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(usr) == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	return usr[0], nil
}

func (p *postgresUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT
			id,
			name,
			email,
			password,
			bio,
			image,
			created_at,
			updated_at
		FROM
			users
		WHERE
			email = $1
			AND is_deleted = FALSE`

	usr, err := p.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	if len(usr) == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	return usr[0], nil
}

func (p *postgresUserRepository) Create(ctx context.Context, usr *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (name, email, password, bio, image)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	if err := p.conn.QueryRow(
		ctx,
		query,
		usr.Name,
		usr.Email,
		usr.Password,
		usr.Bio,
		usr.Image,
	).Scan(
		&usr.ID,
		&usr.CreatedAt,
		&usr.UpdatedAt); err != nil {
		span.SetStatus(codes.Error, "failed inserting users")
		span.RecordError(err)
		return nil, err
	}
	return usr, nil
}

// func (p *postgresUserRepository) Update(ctx context.Context, cat *domain.User) (*domain.User, error) {
// 	query := `
// 		UPDATE users
// 		SET
// 			name = $2,
// 			note = $3,
// 			updated_at = NOW()
// 		WHERE
// 			id = $1
// 		RETURNING updated_at`

// 	ctx, span := spanWithQuery(ctx, p.tracer, query)
// 	defer span.End()

// 	row := p.conn.QueryRow(
// 		ctx,
// 		query,
// 		cat.ID,
// 		cat.Name,
// 		cat.Note,
// 	);

// 	if err := row.Scan(&cat.UpdatedAt); err != nil {
// 		span.SetStatus(codes.Error, "failed to update User")
// 		span.RecordError(err)
// 		return nil, err
// 	}

// 	return cat, nil
// }

// func (p *postgresUserRepository) Delete(ctx context.Context, id int64) error {
// 	query := `
// 		UPDATE users
// 		SET
// 			is_deleted = TRUE,
// 			updated_at = NOW()
// 		WHERE
// 			id = $1`

// 	ctx, span := spanWithQuery(ctx, p.tracer, query)
// 	defer span.End()

// 	result , err := p.conn.Exec(ctx, query, id);

// 	if err != nil {
// 		span.SetStatus(codes.Error, "failed to delete User")
// 		span.RecordError(err)
// 		return err
// 	}

// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return domain.ErrNotFound
// 	}

// 	return nil
// }
