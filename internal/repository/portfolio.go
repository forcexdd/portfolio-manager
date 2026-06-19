package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/repository/db"
)

type postgresPortfolioRepo struct {
	q *db.Queries
}

func NewPostgresPortfolioRepository(sqldb *sql.DB) domain.PortfolioRepository {
	return &postgresPortfolioRepo{q: db.New(sqldb)}
}

func (r *postgresPortfolioRepo) CreatePortfolio(ctx context.Context, name string) (domain.Portfolio, error) {
	p, err := r.q.CreatePortfolio(ctx, name)
	if err != nil {
		return domain.Portfolio{}, err
	}
	return domain.Portfolio{ID: p.ID, Name: p.Name}, nil
}

func (r *postgresPortfolioRepo) GetPortfolio(ctx context.Context, id int64) (domain.Portfolio, error) {
	p, err := r.q.GetPortfolio(ctx, id)
	if err != nil {
		return domain.Portfolio{}, err
	}
	return domain.Portfolio{ID: p.ID, Name: p.Name}, nil
}

func (r *postgresPortfolioRepo) GetPortfolios(ctx context.Context) ([]domain.Portfolio, error) {
	ps, err := r.q.GetAllPortfolios(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Portfolio, len(ps))
	for i, p := range ps {
		res[i] = domain.Portfolio{ID: p.ID, Name: p.Name}
	}
	return res, nil
}

func (r *postgresPortfolioRepo) DeletePortfolio(ctx context.Context, id int64) error {
	return r.q.DeletePortfolio(ctx, id)
}

func (r *postgresPortfolioRepo) UpdatePortfolioAssetLotSize(ctx context.Context, portfolioID, assetID int64, lotSize int64) error {
	return r.q.UpsertPortfolioAssetLot(ctx, db.UpsertPortfolioAssetLotParams{
		PortfolioID: portfolioID,
		AssetID:     assetID,
		LotSize:     lotSize,
	})
}

func (r *postgresPortfolioRepo) UpsertPortfolioAsset(ctx context.Context, portfolioID, assetID, quantity int64) error {
	return r.q.UpsertPortfolioAsset(ctx, db.UpsertPortfolioAssetParams{
		PortfolioID: portfolioID,
		AssetID:     assetID,
		Quantity:    quantity,
	})
}

func (r *postgresPortfolioRepo) GetPortfolioAssets(ctx context.Context, portfolioID int64) ([]domain.PortfolioAsset, error) {
	rows, err := r.q.GetPortfolioAssets(ctx, portfolioID)
	if err != nil {
		return nil, err
	}
	res := make([]domain.PortfolioAsset, len(rows))
	for i, r := range rows {
		p, _ := strconv.ParseFloat(r.Price, 64)
		res[i] = domain.PortfolioAsset{
			Asset:    domain.Asset{ID: r.ID, Name: r.Name, Price: p},
			Quantity: r.Quantity,
			LotSize:  r.LotSize,
		}
	}
	return res, nil
}

func (r *postgresPortfolioRepo) DeletePortfolioAsset(ctx context.Context, portfolioID, assetID int64) error {
	return r.q.DeletePortfolioAsset(ctx, db.DeletePortfolioAssetParams{
		PortfolioID: portfolioID,
		AssetID:     assetID,
	})
}
