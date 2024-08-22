package moex_api_client

import (
	"encoding/json"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_models"
	"io"
	"net/http"
	"strconv"
	"time"
)

type MoexApiClient struct {
	BaseUrl string
}

func NewMoexApiClient() *MoexApiClient {
	return &MoexApiClient{BaseUrl: getBaseUrl()}
}

func getBaseUrl() string {
	return "https://iss.moex.com/iss/"
}

func getCurrentTime() string {
	currTime := time.Now()
	return currTime.Format("2006-01-02")
}

func (m *MoexApiClient) GetAllStocks(time string) ([]*moex_models.StockData, error) {
	url := m.BaseUrl + "statistics/engines/stock/currentprices.json?date=" + time + "&start="
	start := 0
	var allData []*moex_models.StockData
	hasData := true

	for hasData {
		newData, err := m.getStock(url, start)
		if err != nil {
			return nil, err
		}

		if len(newData) == 0 {
			hasData = false
		} else {
			allData = append(allData, newData...)
			start += len(newData)
		}
	}

	return allData, nil
}

func (m *MoexApiClient) getStock(url string, start int) ([]*moex_models.StockData, error) {
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
		//if len(stock) < 8 {
		//	continue
		//}

		newStock := &moex_models.StockData{
			TradeDate:      toString(stock[0]),
			BoardID:        toString(stock[1]),
			SecID:          toString(stock[2]),
			TradeTime:      toString(stock[3]),
			CurPrice:       toFloat64(stock[4]),
			LastPrice:      toFloat64(stock[5]),
			LegalClose:     toFloat64(stock[6]),
			TradingSession: int(toFloat64(stock[7])),
		}

		allData = append(allData, newStock)
	}

	return allData, nil
}

func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func toFloat64(value interface{}) float64 {
	if num, ok := value.(float64); ok {
		return num
	}
	return 0
}
