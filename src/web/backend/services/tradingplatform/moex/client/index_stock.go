package client

import (
	"encoding/json"
	"errors"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"io"
	"net/http"
	"strconv"
)

func (m *MoexApiClient) GetAllIndexAssets(time string, index *model.IndexData) ([]*model.IndexAssetsData, error) {
	url := m.BaseUrl + "statistics/engines/stock/markets/index/analytics/" + index.IndexId + ".json?lang=en&date=" + time + "&start="
	start := 0
	var allData []*model.IndexAssetsData
	hasData := true

	for hasData {
		newData, err := m.getIndexAssetsData(url, start)
		if err != nil {
			return nil, err
		}

		if len(newData) == 0 {
			hasData = false
		} else {
			start += len(newData)
			allData = append(allData, newData...)
		}
	}

	return allData, nil
}

func (m *MoexApiClient) getIndexAssetsData(url string, start int) ([]*model.IndexAssetsData, error) {
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

	var currPricesData *model.IndexAnalyticsData
	currPricesData, err = parseIndexAnalyticsDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseIndexAssetsDataFromIndexAnalytics(currPricesData)
}

func parseIndexAnalyticsDataFromJson(body []byte) (*model.IndexAnalyticsData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data model.IndexAnalyticsData
	err = json.Unmarshal(rawData["analytics"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseIndexAssetsDataFromIndexAnalytics(indexAssetsData *model.IndexAnalyticsData) ([]*model.IndexAssetsData, error) {
	var allData []*model.IndexAssetsData

	for _, indexAsset := range indexAssetsData.Data {
		if len(indexAsset) != 7 {
			return nil, errors.New("invalid index asset data")
		}

		newIndexAsset := &model.IndexAssetsData{
			IndexId:        toString(indexAsset[0]),
			TradeDate:      toString(indexAsset[1]),
			Ticker:         toString(indexAsset[2]),
			ShortNames:     toString(indexAsset[3]),
			SecIds:         toString(indexAsset[4]),
			Weight:         toFloat64(indexAsset[5]),
			TradingSession: int(toFloat64(indexAsset[6])),
		}

		allData = append(allData, newIndexAsset)
	}

	return allData, nil
}
