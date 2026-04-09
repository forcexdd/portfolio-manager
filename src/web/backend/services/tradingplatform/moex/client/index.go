package client

import (
	"encoding/json"
	"errors"
	"io"

	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
)

func (m *MoexApiClient) GetAllIndexes() ([]*moexmodels.IndexData, error) {
	url := m.BaseUrl + "statistics/engines/stock/markets/index/analytics.json?lang=" + language

	response, err := sendGETRequest(url)
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

	var analyticsData *moexmodels.AnalyticsData
	analyticsData, err = parseAnalyticsDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseIndexDataFromAnalytics(analyticsData)
}

func parseAnalyticsDataFromJson(body []byte) (*moexmodels.AnalyticsData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data moexmodels.AnalyticsData
	err = json.Unmarshal(rawData["indices"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseIndexDataFromAnalytics(analyticsData *moexmodels.AnalyticsData) ([]*moexmodels.IndexData, error) {
	var allData []*moexmodels.IndexData

	for _, index := range analyticsData.Data {
		if len(index) != 4 {
			return nil, errors.New("invalid index data")
		}

		newIndex := &moexmodels.IndexData{
			IndexID:   toString(index[0]),
			ShortName: toString(index[1]),
			From:      toString(index[2]),
			Till:      toString(index[3]),
		}

		allData = append(allData, newIndex)
	}

	return allData, nil
}
