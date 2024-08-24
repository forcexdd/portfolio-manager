package stock_exchange_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_models"
)

func (m *MoexService) parseLatestStocks(maxDays int) ([]*moex_models.StockData, error) {
	parseTime := getCurrentTime()

	allStocks, err := m.moexApiClient.GetAllStocks(formatTime(parseTime))
	if err != nil {
		return nil, err
	}
	maxDays--

	for len(allStocks) == 0 && maxDays != 0 {
		parseTime = parseTime.AddDate(0, 0, -1)

		allStocks, err = m.moexApiClient.GetAllStocks(formatTime(parseTime))
		if err != nil {
			return nil, err
		}

		maxDays--
	}

	return allStocks, nil
}

func (m *MoexService) createOrUpdateStock(stock *models.Stock) error {
	dbStock, err := m.StockRepository.GetByName(stock.Name)
	if err != nil {
		return err

	} else if dbStock == nil { // Can't find stock in database
		err = m.StockRepository.Create(stock)
		if err != nil {
			return err
		}
	} else { // Found already existing stock in database
		err = m.StockRepository.Update(stock)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeStockByNameFromSlice(stocks []*models.Stock, name string) []*models.Stock {
	for index, stock := range stocks {
		if stock.Name == name {
			return removeElementFromSliceByIndex[*models.Stock](stocks, index)
		}
	}

	return stocks
}

func (m *MoexService) removeOldStocksFromDb(stocks []*models.Stock) error {
	for _, stock := range stocks {
		err := m.StockRepository.Delete(stock)
		if err != nil {
			return err
		}
	}

	return nil
}
