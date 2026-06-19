package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/forcexdd/portfoliomanager/internal/logger"
	port "github.com/forcexdd/portfoliomanager/internal/port"
	"github.com/forcexdd/portfoliomanager/internal/service"
)

var ErrAlreadyParsed = errors.New("date already parsed")

type ParserWorker struct {
	aSvc   service.AssetService
	iSvc   service.IndexService
	sysSvc service.SystemService
	client port.ExchangeClient
	log    logger.Logger
}

func NewParserWorker(aSvc service.AssetService, iSvc service.IndexService, sysSvc service.SystemService, client port.ExchangeClient, log logger.Logger) *ParserWorker {
	return &ParserWorker{aSvc: aSvc, iSvc: iSvc, sysSvc: sysSvc, client: client, log: log}
}

func (p *ParserWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	workerCtx := logger.IntoContext(ctx, p.log.With("worker", "parser"))

	p.process(workerCtx)
	for {
		select {
		case <-workerCtx.Done():
			return
		case <-ticker.C:
			p.process(workerCtx)
		}
	}
}

func (p *ParserWorker) process(ctx context.Context) {
	log := logger.FromContext(ctx, p.log)

	date, assets, err := p.getValidDateData(ctx)
	if err != nil {
		if errors.Is(err, ErrAlreadyParsed) {
			log.Info("Parser skipped", "reason", "latest date already parsed")
		} else {
			log.Error("Parser failed to find valid asset data", "error", err)
		}
		return
	}

	for _, a := range assets {
		if _, err := p.aSvc.UpsertAsset(ctx, a.Ticker, a.Price); err != nil {
			log.Error("Failed to upsert asset", "sec_id", a.Ticker, "error", err)
		}
	}
	log.Info("Parsing date", "date", date)

	indexes, err := p.client.FetchIndexes(ctx)
	if err != nil {
		log.Error("Parser failed to fetch indexes", "error", err)
		return
	}

	for _, idx := range indexes {
		dbIdx, err := p.iSvc.UpsertIndex(ctx, idx.Ticker)
		if err != nil {
			log.Error("Failed to upsert index", "index_id", idx.Ticker, "error", err)
			continue
		}

		idxAssets, err := p.client.FetchIndexAssets(ctx, date, idx.Ticker)
		if err != nil {
			log.Error("Failed to fetch index assets", "index_id", idx.Ticker, "error", err)
			continue
		}
		if len(idxAssets) == 0 {
			log.Warn("No assets in index", "name", idx.Ticker)
			continue
		}

		for _, ia := range idxAssets {
			asset, err := p.aSvc.GetAssetByName(ctx, ia.Ticker)
			if err != nil {
				log.Warn("No asset from index in DB", "name", idx.Ticker)
				continue
			}
			if err := p.iSvc.UpsertIndexAsset(ctx, dbIdx.ID, asset.ID, ia.Weight/100.0); err != nil {
				log.Error("Failed to upsert index asset", "index_id", dbIdx.ID, "asset_id", asset.ID, "error", err)
			}
		}
	}

	_ = p.sysSvc.MarkDateParsed(ctx, date)
	log.Info("Parsing cycle complete", "date", date)
}

func (p *ParserWorker) getValidDateData(ctx context.Context) (string, []port.MarketAsset, error) {
	now := time.Now()
	for i := 1; i <= 15; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")

		parsed, err := p.sysSvc.HasParsedDate(ctx, date)
		if err == nil && parsed {
			return "", nil, ErrAlreadyParsed
		}

		assets, err := p.client.FetchAssets(ctx, date)
		if err == nil && len(assets) > 0 {
			return date, assets, nil
		}
	}
	return "", nil, fmt.Errorf("exhausted historical search limit")
}
