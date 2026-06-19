package service

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
)

type portfolioService struct {
	pRepo domain.PortfolioRepository
	aRepo domain.AssetRepository
	log   logger.Logger
}

func NewPortfolioService(pRepo domain.PortfolioRepository, aRepo domain.AssetRepository, log logger.Logger) PortfolioService {
	return &portfolioService{pRepo: pRepo, aRepo: aRepo, log: log}
}

func (s *portfolioService) CreatePortfolio(ctx context.Context, name string) (domain.Portfolio, error) {
	log := logger.FromContext(ctx, s.log).With("portfolio_name", name)
	log.Info("Creating portfolio")
	return s.pRepo.CreatePortfolio(ctx, name)
}

func (s *portfolioService) GetPortfolios(ctx context.Context) ([]domain.Portfolio, error) {
	log := logger.FromContext(ctx, s.log)
	log.Debug("Retrieving portfolio indexes list")
	return s.pRepo.GetPortfolios(ctx)
}

func (s *portfolioService) DeletePortfolio(ctx context.Context, id int64) error {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", id)
	log.Info("Deleting portfolio")
	return s.pRepo.DeletePortfolio(ctx, id)
}

func (s *portfolioService) UpdatePortfolioAssetLotSize(ctx context.Context, portfolioID, assetID int64, lotSize int64) error {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", portfolioID, "asset_id", assetID, "lot_size", lotSize)
	log.Info("Updating portfolio asset lot size")
	return s.pRepo.UpdatePortfolioAssetLotSize(ctx, portfolioID, assetID, lotSize)
}

func (s *portfolioService) UpsertPortfolioAsset(ctx context.Context, portfolioID, assetID, quantity int64) error {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", portfolioID, "asset_id", assetID, "quantity", quantity)
	log.Info("Upserting portfolio asset")
	return s.pRepo.UpsertPortfolioAsset(ctx, portfolioID, assetID, quantity)
}

func (s *portfolioService) GetPortfolioAssets(ctx context.Context, portfolioID int64) ([]domain.PortfolioAsset, error) {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", portfolioID)
	log.Debug("Retrieving portfolio assets")
	return s.pRepo.GetPortfolioAssets(ctx, portfolioID)
}

func (s *portfolioService) DeletePortfolioAsset(ctx context.Context, portfolioID, assetID int64) error {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", portfolioID, "asset_id", assetID)
	log.Info("Deleting portfolio asset")
	return s.pRepo.DeletePortfolioAsset(ctx, portfolioID, assetID)
}

func (s *portfolioService) ExportPortfolio(ctx context.Context, id int64) (domain.PortfolioProfile, error) {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", id)
	log.Info("Exporting portfolio layout profile")
	p, err := s.pRepo.GetPortfolio(ctx, id)
	if err != nil {
		return domain.PortfolioProfile{}, err
	}
	assets, err := s.pRepo.GetPortfolioAssets(ctx, id)
	if err != nil {
		return domain.PortfolioProfile{}, err
	}
	profile := domain.PortfolioProfile{Name: p.Name, Assets: []domain.ImportExportAsset{}}
	for _, a := range assets {
		profile.Assets = append(profile.Assets, domain.ImportExportAsset{
			Name:     a.Asset.Name,
			Quantity: a.Quantity,
			LotSize:  a.LotSize,
		})
	}
	return profile, nil
}

func (s *portfolioService) ImportPortfolio(ctx context.Context, name string, assets []domain.ImportExportAsset) (domain.Portfolio, error) {
	log := logger.FromContext(ctx, s.log)
	log.Info("Importing portfolio configuration profile", "name", name, "imported_assets_count", len(assets))
	p, err := s.pRepo.CreatePortfolio(ctx, name)
	if err != nil {
		return domain.Portfolio{}, err
	}
	for _, a := range assets {
		dbAsset, err := s.aRepo.GetAssetByName(ctx, a.Name)
		if err != nil {
			log.Warn("Skipped asset import because it was missing in the global database", "asset_name", a.Name)
			continue
		}
		_ = s.pRepo.UpsertPortfolioAsset(ctx, p.ID, dbAsset.ID, a.Quantity)
		_ = s.pRepo.UpdatePortfolioAssetLotSize(ctx, p.ID, dbAsset.ID, a.LotSize)
	}
	return p, nil
}
