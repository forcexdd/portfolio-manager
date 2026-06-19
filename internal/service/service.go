package service

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
)

type PortfolioService interface {
	CreatePortfolio(ctx context.Context, name string) (domain.Portfolio, error)
	GetPortfolios(ctx context.Context) ([]domain.Portfolio, error)
	DeletePortfolio(ctx context.Context, id int64) error
	UpdatePortfolioAssetLotSize(ctx context.Context, portfolioID, assetID int64, lotSize int64) error
	UpsertPortfolioAsset(ctx context.Context, portfolioID, assetID, quantity int64) error
	GetPortfolioAssets(ctx context.Context, portfolioID int64) ([]domain.PortfolioAsset, error)
	DeletePortfolioAsset(ctx context.Context, portfolioID, assetID int64) error
	ExportPortfolio(ctx context.Context, id int64) (domain.PortfolioProfile, error)
	ImportPortfolio(ctx context.Context, name string, assets []domain.ImportExportAsset) (domain.Portfolio, error)
}

type AssetService interface {
	UpsertAsset(ctx context.Context, name string, price float64) (domain.Asset, error)
	GetAssets(ctx context.Context) ([]domain.Asset, error)
	GetAssetByName(ctx context.Context, name string) (domain.Asset, error)
}

type IndexService interface {
	UpsertIndex(ctx context.Context, name string) (domain.Index, error)
	GetIndexes(ctx context.Context) ([]domain.Index, error)
	UpsertIndexAsset(ctx context.Context, indexID, assetID int64, fraction float64) error
}

type AnalysisService interface {
	AnalyzePortfolio(ctx context.Context, portfolioID, indexID int64) (domain.AnalysisResult, error)
}

type SystemService interface {
	HasParsedDate(ctx context.Context, date string) (bool, error)
	MarkDateParsed(ctx context.Context, date string) error
}
