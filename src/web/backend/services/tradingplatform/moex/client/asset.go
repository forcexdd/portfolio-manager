package client

import (
	"encoding/json"
	"errors"
	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"io"
	"net/http"
	"strconv"
)

// Maybe it's better idea to return set map[*model.AssetData]struct{}
func (m *MoexApiClient) GetAllAssets(time string) ([]*moexmodels.AssetData, error) {
	url := m.BaseUrl + "statistics/engines/stock/currentprices.json?date=" + time + "&start="
	start := 0
	var allData []*moexmodels.AssetData
	hasData := true

	for hasData {
		newData, err := m.getAssetsData(url, start)
		if err != nil {
			return nil, err
		}

		if len(newData) == 0 {
			hasData = false
		} else {
			start += len(newData)
			allData = appendLatestAssetData(allData, newData)
		}
	}

	return allData, nil
}

func (m *MoexApiClient) getAssetsData(url string, start int) ([]*moexmodels.AssetData, error) {
	response, err := http.Get(url + strconv.Itoa(start))
	if err != nil {
		return nil, err
	}

	defer func() {
		newErr := response.Body.Close()
		if newErr != nil {
			err = newErr
		}
	}()

	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var currPricesData *moexmodels.CurrentPricesData
	currPricesData, err = parseCurrentPricesDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseAssetDataFromCurrentPrices(currPricesData)
}

func parseCurrentPricesDataFromJson(body []byte) (*moexmodels.CurrentPricesData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data moexmodels.CurrentPricesData
	err = json.Unmarshal(rawData["currentprices"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseAssetDataFromCurrentPrices(assetData *moexmodels.CurrentPricesData) ([]*moexmodels.AssetData, error) {
	var allData []*moexmodels.AssetData

	for _, asset := range assetData.Data {
		if len(asset) != 8 {
			return nil, errors.New("invalid asset data")
		}

		newAsset := &moexmodels.AssetData{
			TradeDate:      toString(asset[0]),
			BoardID:        toString(asset[1]),
			SecID:          toString(asset[2]),
			TradeTime:      toString(asset[3]),
			CurPrice:       toFloat64(asset[4]),
			LastPrice:      toFloat64(asset[5]),
			LegalClose:     toFloat64(asset[6]),
			TradingSession: int(toFloat64(asset[7])),
		}

		allData = append(allData, newAsset)
	}

	return allData, nil
}

func appendLatestAssetData(allData []*moexmodels.AssetData, newData []*moexmodels.AssetData) []*moexmodels.AssetData {
	latestAssetsMap := make(map[string]*moexmodels.AssetData)
	for _, asset := range newData {
		latestAssetsMap[asset.SecID] = asset
	}

	var updatedLatestAssets []*moexmodels.AssetData
	for i, asset := range allData {
		newLatestAsset, isLatestAsset := latestAssetsMap[asset.SecID] // If we find asset from newData in allData, then update it
		if isLatestAsset {
			allData[i] = newLatestAsset
			updatedLatestAssets = append(updatedLatestAssets, newLatestAsset)
		}
	}

	for _, asset := range latestAssetsMap {
		if !isInside[*moexmodels.AssetData](asset, updatedLatestAssets) {
			allData = append(allData, asset)
		}
	}

	return allData
}
