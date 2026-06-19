package service

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
)

type indexService struct {
	repo domain.IndexRepository
	log  logger.Logger
}

func NewIndexService(repo domain.IndexRepository, log logger.Logger) IndexService {
	return &indexService{repo: repo, log: log}
}

func (s *indexService) UpsertIndex(ctx context.Context, name string) (domain.Index, error) {
	log := logger.FromContext(ctx, s.log).With("index_name", name)
	log.Info("Upserting reference index")
	return s.repo.UpsertIndex(ctx, name)
}

func (s *indexService) GetIndexes(ctx context.Context) ([]domain.Index, error) {
	log := logger.FromContext(ctx, s.log)
	log.Debug("Retrieving reference indexes list")
	return s.repo.GetIndexes(ctx)
}

func (s *indexService) UpsertIndexAsset(ctx context.Context, indexID, assetID int64, fraction float64) error {
	log := logger.FromContext(ctx, s.log).With("index_id", indexID, "asset_id", assetID, "fraction", fraction)
	log.Debug("Upserting index reference asset")
	return s.repo.UpsertIndexAsset(ctx, indexID, assetID, fraction)
}
