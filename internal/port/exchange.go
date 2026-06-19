package port

import "context"

type MarketAsset struct {
	Ticker string
	Price  float64
}

type MarketIndex struct {
	Ticker string
}

type MarketIndexAsset struct {
	Ticker string
	Weight float64
}

type ExchangeClient interface {
	FetchAssets(ctx context.Context, date string) ([]MarketAsset, error)
	FetchIndexes(ctx context.Context) ([]MarketIndex, error)
	FetchIndexAssets(ctx context.Context, date, indexID string) ([]MarketIndexAsset, error)
}
