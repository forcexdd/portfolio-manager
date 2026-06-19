package repository

import (
	"context"
	"database/sql"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/repository/db"
)

type postgresSystemRepo struct {
	q *db.Queries
}

func NewPostgresSystemRepository(sqldb *sql.DB) domain.SystemRepository {
	return &postgresSystemRepo{q: db.New(sqldb)}
}

func (r *postgresSystemRepo) HasParsedDate(ctx context.Context, date string) (bool, error) {
	return r.q.HasParsedDate(ctx, date)
}

func (r *postgresSystemRepo) MarkDateParsed(ctx context.Context, date string) error {
	return r.q.MarkDateParsed(ctx, date)
}
