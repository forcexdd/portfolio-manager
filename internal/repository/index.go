package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/repository/db"
)

type indexRepo struct {
	q *db.Queries
}

func NewPostgresIndexRepository(sqldb *sql.DB) domain.IndexRepository {
	return &indexRepo{q: db.New(sqldb)}
}

func (r *indexRepo) UpsertIndex(ctx context.Context, name string) (domain.Index, error) {
	i, err := r.q.UpsertIndex(ctx, name)
	if err != nil {
		return domain.Index{}, err
	}
	return domain.Index{ID: i.ID, Name: i.Name}, nil
}

func (r *indexRepo) GetIndexes(ctx context.Context) ([]domain.Index, error) {
	is, err := r.q.GetAllIndexes(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Index, len(is))
	for i, idx := range is {
		res[i] = domain.Index{ID: idx.ID, Name: idx.Name}
	}
	return res, nil
}

func (r *indexRepo) UpsertIndexAsset(ctx context.Context, indexID, assetID int64, fraction float64) error {
	fStr := strconv.FormatFloat(fraction, 'f', -1, 64)
	return r.q.UpsertIndexAsset(ctx, db.UpsertIndexAssetParams{
		IndexID:  indexID,
		AssetID:  assetID,
		Fraction: fStr,
	})
}

func (r *indexRepo) GetIndexAssets(ctx context.Context, indexID int64) ([]domain.IndexAsset, error) {
	rows, err := r.q.GetIndexAssets(ctx, indexID)
	if err != nil {
		return nil, err
	}
	res := make([]domain.IndexAsset, len(rows))
	for i, r := range rows {
		p, _ := strconv.ParseFloat(r.Price, 64)
		f, _ := strconv.ParseFloat(r.Fraction, 64)
		res[i] = domain.IndexAsset{
			Asset:    domain.Asset{ID: r.ID, Name: r.Name, Price: p},
			Fraction: f,
		}
	}
	return res, nil
}
