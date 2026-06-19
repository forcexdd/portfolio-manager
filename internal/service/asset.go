package service

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
)

type assetService struct {
	repo domain.AssetRepository
	log  logger.Logger
}

func NewAssetService(repo domain.AssetRepository, log logger.Logger) AssetService {
	return &assetService{repo: repo, log: log}
}

func (s *assetService) UpsertAsset(ctx context.Context, name string, price float64) (domain.Asset, error) {
	log := logger.FromContext(ctx, s.log).With("asset_name", name, "price", price)
	log.Debug("Upserting system asset")
	return s.repo.UpsertAsset(ctx, name, price)
}

func (s *assetService) GetAssets(ctx context.Context) ([]domain.Asset, error) {
	log := logger.FromContext(ctx, s.log)
	log.Debug("Retrieving system assets")
	return s.repo.GetAssets(ctx)
}

func (s *assetService) GetAssetByName(ctx context.Context, name string) (domain.Asset, error) {
	log := logger.FromContext(ctx, s.log).With("asset_name", name)
	log.Debug("Fetching system asset by name")
	return s.repo.GetAssetByName(ctx, name)
}
