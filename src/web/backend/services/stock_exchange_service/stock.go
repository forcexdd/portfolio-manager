package stock_exchange_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
)

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
