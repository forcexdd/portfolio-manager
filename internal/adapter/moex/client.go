package moex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/forcexdd/portfoliomanager/internal/port"
)

type MoexClient struct {
	baseURL string
	client  *http.Client
}

func NewMoexClient(baseURL string, client *http.Client) *MoexClient {
	return &MoexClient{baseURL: baseURL, client: client}
}

func (m *MoexClient) FetchAssets(ctx context.Context, date string) ([]port.MarketAsset, error) {
	var assets []port.MarketAsset
	start := 0

	for {
		url := fmt.Sprintf("%s/statistics/engines/stock/currentprices.json?date=%s&start=%d", m.baseURL, date, start)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := m.client.Do(req)
		if err != nil {
			return nil, err
		}

		var result struct {
			CurrentPrices struct {
				Data [][]interface{} `json:"data"`
			} `json:"currentprices"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if len(result.CurrentPrices.Data) == 0 {
			break
		}

		pageAssetsCount := 0
		for _, row := range result.CurrentPrices.Data {
			if len(row) > 6 && row[3] != nil && row[5] != nil {
				if secID, ok1 := row[3].(string); ok1 {
					if price, ok2 := row[5].(float64); ok2 {
						assets = append(assets, port.MarketAsset{Ticker: secID, Price: price})
						pageAssetsCount++
					}
				}
			}
		}

		if pageAssetsCount == 0 {
			break
		}
		start += len(result.CurrentPrices.Data)
	}

	return assets, nil
}

func (m *MoexClient) FetchIndexes(ctx context.Context) ([]port.MarketIndex, error) {
	url := fmt.Sprintf("%s/statistics/engines/stock/markets/index/analytics.json", m.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Indices struct {
			Data [][]interface{} `json:"data"`
		} `json:"indices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var indexes []port.MarketIndex
	for _, row := range result.Indices.Data {
		if len(row) > 0 && row[0] != nil {
			if id, ok := row[0].(string); ok {
				indexes = append(indexes, port.MarketIndex{Ticker: id})
			}
		}
	}
	return indexes, nil
}

func (m *MoexClient) FetchIndexAssets(ctx context.Context, date, indexID string) ([]port.MarketIndexAsset, error) {
	var assets []port.MarketIndexAsset
	start := 0

	for {
		url := fmt.Sprintf("%s/statistics/engines/stock/markets/index/analytics/%s.json?date=%s&start=%d", m.baseURL, indexID, date, start)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := m.client.Do(req)
		if err != nil {
			return nil, err
		}

		var result struct {
			Analytics struct {
				Data [][]interface{} `json:"data"`
			} `json:"analytics"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if len(result.Analytics.Data) == 0 {
			break
		}

		pageAssetsCount := 0
		for _, row := range result.Analytics.Data {
			if len(row) > 5 && row[4] != nil && row[5] != nil {
				if sec, ok1 := row[4].(string); ok1 {
					if weight, ok2 := row[5].(float64); ok2 {
						assets = append(assets, port.MarketIndexAsset{Ticker: sec, Weight: weight})
						pageAssetsCount++
					}
				}
			}
		}

		if pageAssetsCount == 0 {
			break
		}
		start += len(result.Analytics.Data)
	}

	return assets, nil
}
