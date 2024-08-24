package stock_exchange_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_models"
)

func (m *MoexService) createOrUpdateIndex(index *models.Index) error {
	dbIndex, err := m.IndexRepository.GetByName(index.Name)
	if err != nil {
		return err

	} else if dbIndex == nil { // Can't find index in database
		err = m.IndexRepository.Create(index)
		if err != nil {
			return err
		}
	} else { // Found already existing index in database
		err = m.IndexRepository.Update(index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) createStocksFractionMapFromIndexStocks(indexStocks []*moex_models.IndexStocksData) (map[*models.Stock]float64, error) {
	newStocksFractionMap := make(map[*models.Stock]float64)
	for _, indexStock := range indexStocks {
		stock, err := m.StockRepository.GetByName(indexStock.SecIds)
		if err != nil {
			return nil, err
		}

		newStocksFractionMap[stock] = indexStock.Weight
	}

	return newStocksFractionMap, nil
}

func removeIndexByNameFromSlice(indexes []*models.Index, name string) []*models.Index {
	for i, index := range indexes {
		if index.Name == name {
			return removeElementFromSliceByIndex[*models.Index](indexes, i)
		}
	}

	return indexes
}

func (m *MoexService) removeOldIndexesFromDb(indexes []*models.Index) error {
	for _, index := range indexes {
		err := m.IndexRepository.Delete(index)
		if err != nil {
			return err
		}
	}

	return nil
}
