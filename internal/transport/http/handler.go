package http

import (
	"context"

	"github.com/forcexdd/portfoliomanager/internal/domain"
	"github.com/forcexdd/portfoliomanager/internal/logger"
	"github.com/forcexdd/portfoliomanager/internal/service"
	"github.com/forcexdd/portfoliomanager/internal/transport/http/api"
)

type Handler struct {
	pSvc  service.PortfolioService
	aSvc  service.AssetService
	iSvc  service.IndexService
	anSvc service.AnalysisService
	log   logger.Logger
}

func NewHandler(pSvc service.PortfolioService, aSvc service.AssetService, iSvc service.IndexService, anSvc service.AnalysisService, log logger.Logger) api.Handler {
	return &Handler{pSvc: pSvc, aSvc: aSvc, iSvc: iSvc, anSvc: anSvc, log: log}
}

func (h *Handler) ListPortfolios(ctx context.Context) ([]api.Portfolio, error) {
	ports, err := h.pSvc.GetPortfolios(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]api.Portfolio, len(ports))
	for i, p := range ports {
		res[i] = api.Portfolio{ID: p.ID, Name: p.Name}
	}
	return res, nil
}

func (h *Handler) CreatePortfolio(ctx context.Context, req *api.CreatePortfolioReq) (*api.Portfolio, error) {
	p, err := h.pSvc.CreatePortfolio(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &api.Portfolio{ID: p.ID, Name: p.Name}, nil
}

func (h *Handler) DeletePortfolio(ctx context.Context, params api.DeletePortfolioParams) error {
	return h.pSvc.DeletePortfolio(ctx, params.ID)
}

func (h *Handler) ListPortfolioAssets(ctx context.Context, params api.ListPortfolioAssetsParams) ([]api.PortfolioAsset, error) {
	assets, err := h.pSvc.GetPortfolioAssets(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	res := make([]api.PortfolioAsset, len(assets))
	for i, a := range assets {
		res[i] = api.PortfolioAsset{
			Asset: api.Asset{
				ID:    a.Asset.ID,
				Name:  a.Asset.Name,
				Price: a.Asset.Price,
			},
			Quantity: a.Quantity,
			LotSize:  a.LotSize,
		}
	}
	return res, nil
}

func (h *Handler) UpsertPortfolioAsset(ctx context.Context, req *api.UpsertPortfolioAssetReq, params api.UpsertPortfolioAssetParams) error {
	return h.pSvc.UpsertPortfolioAsset(ctx, params.ID, req.AssetID, req.Quantity)
}

func (h *Handler) DeletePortfolioAsset(ctx context.Context, params api.DeletePortfolioAssetParams) error {
	return h.pSvc.DeletePortfolioAsset(ctx, params.ID, params.AssetID)
}

func (h *Handler) UpdatePortfolioAssetLotSize(ctx context.Context, req *api.UpdateLotSizeReq, params api.UpdatePortfolioAssetLotSizeParams) error {
	return h.pSvc.UpdatePortfolioAssetLotSize(ctx, params.ID, params.AssetID, req.LotSize)
}

func (h *Handler) ListAssets(ctx context.Context) ([]api.Asset, error) {
	assets, err := h.aSvc.GetAssets(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]api.Asset, len(assets))
	for i, a := range assets {
		res[i] = api.Asset{ID: a.ID, Name: a.Name, Price: a.Price}
	}
	return res, nil
}

func (h *Handler) ListIndexes(ctx context.Context) ([]api.Index, error) {
	indexes, err := h.iSvc.GetIndexes(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]api.Index, len(indexes))
	for i, idx := range indexes {
		res[i] = api.Index{ID: idx.ID, Name: idx.Name}
	}
	return res, nil
}

func (h *Handler) AnalyzePortfolio(ctx context.Context, params api.AnalyzePortfolioParams) (*api.AnalysisResult, error) {
	res, err := h.anSvc.AnalyzePortfolio(ctx, params.PortfolioID, params.IndexID)
	if err != nil {
		return nil, err
	}

	apiAssets := make([]api.AssetDiff, len(res.Assets))
	for i, a := range res.Assets {
		apiAssets[i] = api.AssetDiff{
			AssetID:         a.AssetID,
			Name:            a.Name,
			Price:           a.Price,
			CurrentQuantity: a.CurrentQuantity,
			CurrentFraction: a.CurrentFraction,
			IndexFraction:   a.IndexFraction,
			FractionDiff:    a.FractionDiff,
			TargetQuantity:  a.TargetQuantity,
			Difference:      a.Difference,
			LotSize:         a.LotSize,
			Action:          a.Action,
		}
	}
	return &api.AnalysisResult{
		TotalValue: res.TotalValue,
		Assets:     apiAssets,
	}, nil
}

func (h *Handler) ExportPortfolio(ctx context.Context, params api.ExportPortfolioParams) (*api.ExportPortfolioResult, error) {
	res, err := h.pSvc.ExportPortfolio(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	assets := make([]api.ImportExportAsset, len(res.Assets))
	for i, a := range res.Assets {
		assets[i] = api.ImportExportAsset{
			Name:     a.Name,
			Quantity: a.Quantity,
			LotSize:  a.LotSize,
		}
	}
	return &api.ExportPortfolioResult{
		Name:   res.Name,
		Assets: assets,
	}, nil
}

func (h *Handler) ImportPortfolio(ctx context.Context, req *api.ImportPortfolioReq) (*api.Portfolio, error) {
	assets := make([]domain.ImportExportAsset, len(req.Assets))
	for i, a := range req.Assets {
		assets[i] = domain.ImportExportAsset{
			Name:     a.Name,
			Quantity: a.Quantity,
			LotSize:  a.LotSize,
		}
	}
	p, err := h.pSvc.ImportPortfolio(ctx, req.Name, assets)
	if err != nil {
		return nil, err
	}
	return &api.Portfolio{ID: p.ID, Name: p.Name}, nil
}
