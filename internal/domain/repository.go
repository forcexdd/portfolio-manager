package domain

import "context"

type PortfolioRepository interface {
	CreatePortfolio(ctx context.Context, name string) (Portfolio, error)
	GetPortfolio(ctx context.Context, id int64) (Portfolio, error)
	GetPortfolios(ctx context.Context) ([]Portfolio, error)
	DeletePortfolio(ctx context.Context, id int64) error

	UpsertPortfolioAsset(ctx context.Context, portfolioID, assetID, quantity int64) error
	UpdatePortfolioAssetLotSize(ctx context.Context, portfolioID, assetID int64, lotSize int64) error
	GetPortfolioAssets(ctx context.Context, portfolioID int64) ([]PortfolioAsset, error)
	DeletePortfolioAsset(ctx context.Context, portfolioID, assetID int64) error
}

type AssetRepository interface {
	UpsertAsset(ctx context.Context, name string, price float64) (Asset, error)
	GetAssets(ctx context.Context) ([]Asset, error)
	GetAssetByName(ctx context.Context, name string) (Asset, error)
}

type IndexRepository interface {
	UpsertIndex(ctx context.Context, name string) (Index, error)
	GetIndexes(ctx context.Context) ([]Index, error)
	UpsertIndexAsset(ctx context.Context, indexID, assetID int64, fraction float64) error
	GetIndexAssets(ctx context.Context, indexID int64) ([]IndexAsset, error)
}

type SystemRepository interface {
	HasParsedDate(ctx context.Context, date string) (bool, error)
	MarkDateParsed(ctx context.Context, date string) error
}
