package service

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
)

type systemService struct {
	repo domain.SystemRepository
	log  logger.Logger
}

func NewSystemService(repo domain.SystemRepository, log logger.Logger) SystemService {
	return &systemService{repo: repo, log: log}
}

func (s *systemService) HasParsedDate(ctx context.Context, date string) (bool, error) {
	log := logger.FromContext(ctx, s.log).With("date", date)
	log.Debug("Checking database parsing records")
	return s.repo.HasParsedDate(ctx, date)
}

func (s *systemService) MarkDateParsed(ctx context.Context, date string) error {
	log := logger.FromContext(ctx, s.log).With("date", date)
	log.Info("Marking system state date as successfully processed")
	return s.repo.MarkDateParsed(ctx, date)
}
