package moex_api_client

import (
	"encoding/json"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_models"
	"io"
	"net/http"
)

func (m *MoexApiClient) GetAllIndexes() ([]*moex_models.IndexData, error) {
	url := m.BaseUrl + "statistics/engines/stock/markets/index/analytics.json?lang=en"

	response, err := http.Get(url)
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

	var analyticsData *moex_models.AnalyticsData
	analyticsData, err = parseAnalyticsDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseIndexDataFromAnalytics(analyticsData)
}

func parseAnalyticsDataFromJson(body []byte) (*moex_models.AnalyticsData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data moex_models.AnalyticsData
	err = json.Unmarshal(rawData["indices"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseIndexDataFromAnalytics(analyticsData *moex_models.AnalyticsData) ([]*moex_models.IndexData, error) {
	var allData []*moex_models.IndexData

	for _, index := range analyticsData.Data {
		if len(index) != 4 {
			return nil, errors.New("invalid index data")
		}

		newIndex := &moex_models.IndexData{
			IndexId:   toString(index[0]),
			ShortName: toString(index[1]),
			From:      toString(index[2]),
			Till:      toString(index[3]),
		}

		allData = append(allData, newIndex)
	}

	return allData, nil
}
