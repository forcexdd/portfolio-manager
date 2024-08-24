package moex_api_client

import (
	"encoding/json"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_models"
	"io"
	"net/http"
	"strconv"
)

func (m *MoexApiClient) GetAllIndexStocks(time string, index *moex_models.IndexData) ([]*moex_models.IndexStocksData, error) {
	url := m.BaseUrl + "statistics/engines/stock/markets/index/analytics/" + index.IndexId + ".json?lang=en&date=" + time + "&start="
	start := 0
	var allData []*moex_models.IndexStocksData
	hasData := true

	for hasData {
		newData, err := m.getIndexStocksData(url, start)
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

func (m *MoexApiClient) getIndexStocksData(url string, start int) ([]*moex_models.IndexStocksData, error) {
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

	var currPricesData *moex_models.IndexAnalyticsData
	currPricesData, err = parseIndexAnalyticsDataFromJson(body)
	if err != nil {
		return nil, err
	}

	return parseIndexStocksDataFromIndexAnalytics(currPricesData)
}

func parseIndexAnalyticsDataFromJson(body []byte) (*moex_models.IndexAnalyticsData, error) {
	var rawData map[string]json.RawMessage
	err := json.Unmarshal(body, &rawData)
	if err != nil {
		return nil, err
	}

	var data moex_models.IndexAnalyticsData
	err = json.Unmarshal(rawData["analytics"], &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func parseIndexStocksDataFromIndexAnalytics(indexStocksData *moex_models.IndexAnalyticsData) ([]*moex_models.IndexStocksData, error) {
	var allData []*moex_models.IndexStocksData

	for _, indexStock := range indexStocksData.Data {
		if len(indexStock) != 7 {
			return nil, errors.New("invalid index stock data")
		}

		newIndexStock := &moex_models.IndexStocksData{
			IndexId:        toString(indexStock[0]),
			TradeDate:      toString(indexStock[1]),
			Ticker:         toString(indexStock[2]),
			ShortNames:     toString(indexStock[3]),
			SecIds:         toString(indexStock[4]),
			Weight:         toFloat64(indexStock[5]),
			TradingSession: int(toFloat64(indexStock[6])),
		}

		if len(newIndexStock.SecIds) > 5 {
			return nil, nil // Index contains bonds
		}

		allData = append(allData, newIndexStock)
	}

	return allData, nil
}
