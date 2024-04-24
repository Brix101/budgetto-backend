package util

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Brix101/budgetto-backend/internal/repository"
)

func spanWithQuery(ctx context.Context, tracer trace.Tracer, query string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "db:query")
	span.SetAttributes(semconv.DBStatementKey.String(query))
	return ctx, span
}

type postgresSeedRepository struct {
	ctx    context.Context
	conn   repository.Connection
	tracer trace.Tracer
	logger *zap.Logger
}

func NewSeeder(ctx context.Context, logger *zap.Logger, conn repository.Connection) *postgresSeedRepository {
	tracer := otel.Tracer("db:postgres:seeder")

	return &postgresSeedRepository{ctx: ctx, conn: conn, logger: logger, tracer: tracer}
}

func (p *postgresSeedRepository) CategorySeed() error {
	ctx := p.ctx
	logger := p.logger
	t_query := `
		SELECT COUNT(*) FROM categories`

	t_ctx, t_span := spanWithQuery(ctx, p.tracer, t_query)
	defer t_span.End()

	result := p.conn.QueryRow(t_ctx, t_query)

	var count int
	result.Scan(&count)

	if count <= 0 {
		query := `
			INSERT INTO categories (name, note)
			VALUES 
		        ('Debt Payments', 'This category can include payments towards credit card debt, student loans, or other debts.'),
                ('Entertainment', 'This category can include expenses for movies, concerts, hobbies, and vacations.'),
                ('Food', 'This category can include groceries, dining out, and snacks.'),
                ('Health Care', 'This category can include expenses for health insurance, doctor visits, prescriptions, and other medical expenses.'),
                ('Housing', 'This category can include mortgage or rent payments, property taxes, homeowners or renters insurance, repairs and maintenance, and utilities.'),
                ('Personal Care', 'This category can include items such as haircuts, personal grooming products, and gym memberships.'),
                ('Savings', 'This category can include savings towards retirement, emergency funds, or other financial goals.'),
                ('Transportation', 'This category can include car payments, gas, car insurance, maintenance and repairs, and public transportation expenses.'),
                ('Utilities', 'This category can include expenses for electricity, gas, water, internet, and phone.')`

		ctx, span := spanWithQuery(ctx, p.tracer, query)
		defer span.End()
		result, err := p.conn.Exec(ctx, query)
		if err != nil {
			span.SetStatus(codes.Error, "failed to seed category")
			span.RecordError(err)
			logger.Error("❌❌❌ Failed to seed category:", zap.Error(err))
			return err
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected >= 1 {
			logger.Info("✅✅✅ Category seeder executed successfully.")
		}
		return nil
	}

	logger.Info("👍👍👍 Category records already exist. Skipping the seeder.")
	return nil
}
