package stock_exchange_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_api_client"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_models"
	"time"
)

type StockExchangeService interface {
	ParseAllStocksIntoDb() error
	ParseAllIndexesIntoDb() error
}

type MoexService struct {
	StockRepository repositories.StockRepository
	IndexRepository repositories.IndexRepository
	moexApiClient   *moex_api_client.MoexApiClient
}

func NewStockExchangeService(stockRepository repositories.StockRepository, indexRepository repositories.IndexRepository) StockExchangeService {
	newService := &MoexService{
		StockRepository: stockRepository,
		IndexRepository: indexRepository,
	}
	newService.setApiClient()

	return newService
}

func (m *MoexService) setApiClient() {
	m.moexApiClient = moex_api_client.NewMoexApiClient()
}

func (m *MoexService) ParseAllStocksIntoDb() error {
	allStocks, err := m.moexApiClient.GetAllStocks(getCurrentTime())
	if err != nil {
		return err
	}

	for _, stock := range allStocks {
		newStock := &models.Stock{
			Name:  stock.SecId,
			Price: stock.CurPrice,
		}

		err = m.StockRepository.Create(newStock)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) ParseAllIndexesIntoDb() error {
	allIndexes, err := m.moexApiClient.GetAllIndexes()
	if err != nil {
		return err
	}

	for _, index := range allIndexes {
		var indexStocks []*moex_models.IndexStocksData
		indexStocks, err = m.moexApiClient.GetAllIndexStocks(getCurrentTime(), index)
		if err != nil {
			return err
		}
		if indexStocks == nil {
			continue // Index contains bonds
		}

		newStocksFractionMap := make(map[*models.Stock]float64)
		for _, indexStock := range indexStocks {
			var stock *models.Stock
			stock, err = m.StockRepository.GetByName(indexStock.SecIds)
			if err != nil {
				return err
			}

			newStocksFractionMap[stock] = indexStock.Weight
		}

		newIndex := &models.Index{
			Name:              index.IndexId,
			StocksFractionMap: newStocksFractionMap,
		}

		err = m.IndexRepository.Create(newIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCurrentTime() string {
	currTime := time.Now()
	return currTime.Format("2006-01-02")
}
