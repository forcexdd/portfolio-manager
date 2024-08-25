package moex_api_client

import (
	"encoding/json"
	"errors"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_models"
	"io"
	"net/http"
	"strconv"
)

// Maybe it's better idea to return set map[*moex_models.AssetData]struct{}
func (m *MoexApiClient) GetAllAssets(time string) ([]*moex_models.AssetData, error) {
	url := m.BaseUrl + "statistics/engines/stock/currentprices.json?date=" + time + "&start="
	start := 0
	var allData []*moex_models.AssetData
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

func (m *MoexApiClient) getAssetsData(url string, start int) ([]*moex_models.AssetData, error) {
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

	var currPricesData *moex_models.CurrentPricesData
	currPricesData, err = parseCurrentPricesDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseAssetDataFromCurrentPrices(currPricesData)
}

func parseCurrentPricesDataFromJson(body []byte) (*moex_models.CurrentPricesData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data moex_models.CurrentPricesData
	err = json.Unmarshal(rawData["currentprices"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseAssetDataFromCurrentPrices(assetData *moex_models.CurrentPricesData) ([]*moex_models.AssetData, error) {
	var allData []*moex_models.AssetData

	for _, asset := range assetData.Data {
		if len(asset) != 8 {
			return nil, errors.New("invalid asset data")
		}

		newAsset := &moex_models.AssetData{
			TradeDate:      toString(asset[0]),
			BoardId:        toString(asset[1]),
			SecId:          toString(asset[2]),
			TradeTime:      toString(asset[3]),
			CurPrice:       toFloat64(asset[4]),
			LastPrice:      toFloat64(asset[5]),
			LegalClose:     toFloat64(asset[6]),
			TradingSession: int(toFloat64(asset[7])),
		}

		if len(newAsset.SecId) <= 5 { // Since we are getting only stocks. Remove this to get bonds in addition
			allData = append(allData, newAsset)
		}
	}

	return allData, nil
}

func appendLatestAssetData(allData []*moex_models.AssetData, newData []*moex_models.AssetData) []*moex_models.AssetData {
	latestAssetsMap := make(map[string]*moex_models.AssetData)
	for _, asset := range newData {
		latestAssetsMap[asset.SecId] = asset
	}

	var updatedLatestAssets []*moex_models.AssetData
	for i, asset := range allData {
		newLatestAsset, isLatestAsset := latestAssetsMap[asset.SecId] // If we find asset from newData in allData, then update it
		if isLatestAsset {
			allData[i] = newLatestAsset
			updatedLatestAssets = append(updatedLatestAssets, newLatestAsset)
		}
	}

	for _, asset := range latestAssetsMap {
		if !isInside[*moex_models.AssetData](asset, updatedLatestAssets) {
			allData = append(allData, asset)
		}
	}

	return allData
}
