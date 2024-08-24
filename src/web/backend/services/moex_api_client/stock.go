package moex_api_client

import (
	"encoding/json"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_models"
	"io"
	"net/http"
	"strconv"
)

// Maybe it's better idea to return set map[*moex_models.StockData]struct{}
func (m *MoexApiClient) GetAllStocks(time string) ([]*moex_models.StockData, error) {
	url := m.BaseUrl + "statistics/engines/stock/currentprices.json?date=" + time + "&start="
	start := 0
	var allData []*moex_models.StockData
	hasData := true

	for hasData {
		newData, err := m.getStocksData(url, start)
		if err != nil {
			return nil, err
		}

		if len(newData) == 0 {
			hasData = false
		} else {
			start += len(newData)
			allData = appendLatestStockData(allData, newData)
		}
	}

	return allData, nil
}

func (m *MoexApiClient) getStocksData(url string, start int) ([]*moex_models.StockData, error) {
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

	return parseStockDataFromCurrentPrices(currPricesData)
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

func parseStockDataFromCurrentPrices(stockData *moex_models.CurrentPricesData) ([]*moex_models.StockData, error) {
	var allData []*moex_models.StockData

	for _, stock := range stockData.Data {
		if len(stock) != 8 {
			return nil, errors.New("invalid stock data")
		}

		newStock := &moex_models.StockData{
			TradeDate:      toString(stock[0]),
			BoardId:        toString(stock[1]),
			SecId:          toString(stock[2]),
			TradeTime:      toString(stock[3]),
			CurPrice:       toFloat64(stock[4]),
			LastPrice:      toFloat64(stock[5]),
			LegalClose:     toFloat64(stock[6]),
			TradingSession: int(toFloat64(stock[7])),
		}

		if len(newStock.SecId) <= 5 { // Since we are getting only stocks. Remove this to get bonds in addition
			allData = append(allData, newStock)
		}
	}

	return allData, nil
}

func appendLatestStockData(allData []*moex_models.StockData, newData []*moex_models.StockData) []*moex_models.StockData {
	latestStocksMap := make(map[string]*moex_models.StockData)
	for _, stock := range newData {
		latestStocksMap[stock.SecId] = stock
	}

	var updatedLatestStocks []*moex_models.StockData
	for i, stock := range allData {
		newLatestStock, isLatestStock := latestStocksMap[stock.SecId] // If we find stock from newData in allData, then update it
		if isLatestStock {
			allData[i] = newLatestStock
			updatedLatestStocks = append(updatedLatestStocks, newLatestStock)
		}
	}

	for _, stock := range latestStocksMap {
		if !isInside[*moex_models.StockData](stock, updatedLatestStocks) {
			allData = append(allData, stock)
		}
	}

	return allData
}
