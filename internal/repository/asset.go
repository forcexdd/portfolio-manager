package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/repository/db"
)

type postgresAssetRepo struct {
	q *db.Queries
}

func NewPostgresAssetRepository(sqldb *sql.DB) domain.AssetRepository {
	return &postgresAssetRepo{q: db.New(sqldb)}
}

func (r *postgresAssetRepo) UpsertAsset(ctx context.Context, name string, price float64) (domain.Asset, error) {
	priceStr := strconv.FormatFloat(price, 'f', -1, 64)
	a, err := r.q.UpsertAsset(ctx, db.UpsertAssetParams{Name: name, Price: priceStr})
	if err != nil {
		return domain.Asset{}, err
	}
	p, _ := strconv.ParseFloat(a.Price, 64)
	return domain.Asset{ID: a.ID, Name: a.Name, Price: p}, nil
}

func (r *postgresAssetRepo) GetAssets(ctx context.Context) ([]domain.Asset, error) {
	as, err := r.q.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Asset, len(as))
	for i, a := range as {
		p, _ := strconv.ParseFloat(a.Price, 64)
		res[i] = domain.Asset{ID: a.ID, Name: a.Name, Price: p}
	}
	return res, nil
}

func (r *postgresAssetRepo) GetAssetByName(ctx context.Context, name string) (domain.Asset, error) {
	a, err := r.q.GetAssetByName(ctx, name)
	if err != nil {
		return domain.Asset{}, err
	}
	p, _ := strconv.ParseFloat(a.Price, 64)
	return domain.Asset{ID: a.ID, Name: a.Name, Price: p}, nil
}
