package service

import (
	"context"
	"math"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
)

type analysisService struct {
	pRepo domain.PortfolioRepository
	iRepo domain.IndexRepository
	log   logger.Logger
}

func NewAnalysisService(pRepo domain.PortfolioRepository, iRepo domain.IndexRepository, log logger.Logger) AnalysisService {
	return &analysisService{pRepo: pRepo, iRepo: iRepo, log: log}
}

func (s *analysisService) AnalyzePortfolio(ctx context.Context, portfolioID, indexID int64) (domain.AnalysisResult, error) {
	log := logger.FromContext(ctx, s.log).With("portfolio_id", portfolioID, "index_id", indexID)
	log.Info("Performing target alignment tracking analysis")

	pAssets, err := s.pRepo.GetPortfolioAssets(ctx, portfolioID)
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	iAssets, err := s.iRepo.GetIndexAssets(ctx, indexID)
	if err != nil {
		return domain.AnalysisResult{}, err
	}

	var totalValue float64
	currentMap := make(map[int64]domain.PortfolioAsset)
	for _, pa := range pAssets {
		totalValue += pa.Asset.Price * float64(pa.Quantity)
		currentMap[pa.Asset.ID] = pa
	}

	targetMap := make(map[int64]domain.IndexAsset)
	for _, ia := range iAssets {
		targetMap[ia.Asset.ID] = ia
	}

	allIDs := make(map[int64]bool)
	for id := range currentMap {
		allIDs[id] = true
	}
	for id := range targetMap {
		allIDs[id] = true
	}

	var res domain.AnalysisResult
	res.TotalValue = totalValue

	for id := range allIDs {
		c, cOk := currentMap[id]
		t, tOk := targetMap[id]

		asset := c.Asset
		if !cOk {
			asset = t.Asset
		}

		lotSize := c.LotSize
		if lotSize <= 0 {
			lotSize = 1
		}

		var currentQuantity int64
		if cOk {
			currentQuantity = c.Quantity
		}

		var indexFraction float64
		if tOk {
			indexFraction = t.Fraction
		}

		currentFraction := 0.0
		if totalValue > 0 {
			currentFraction = (asset.Price * float64(currentQuantity)) / totalValue
		}

		targetValue := totalValue * indexFraction
		rawTargetQty := targetValue / asset.Price

		lotSizeFloat := float64(lotSize)
		roundedTargetQty := math.Round(rawTargetQty/lotSizeFloat) * lotSizeFloat
		targetQuantity := int64(roundedTargetQty)

		diffQty := targetQuantity - currentQuantity
		action := "HOLD"
		if diffQty > 0 {
			action = "BUY"
		} else if diffQty < 0 {
			action = "SELL"
		}

		res.Assets = append(res.Assets, domain.AssetDiff{
			AssetID:         asset.ID,
			Name:            asset.Name,
			Price:           asset.Price,
			CurrentQuantity: currentQuantity,
			CurrentFraction: currentFraction * 100,
			IndexFraction:   indexFraction * 100,
			FractionDiff:    (currentFraction - indexFraction) * 100,
			TargetQuantity:  targetQuantity,
			Difference:      diffQty,
			LotSize:         lotSize,
			Action:          action,
		})
	}

	return res, nil
}
