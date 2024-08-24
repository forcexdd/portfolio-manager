package stock_exchange_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_api_client"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service/moex/moex_models"
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

func (m *MoexService) ParseAllStocksIntoDb() error {
	allStocks, err := m.parseLatestStocks(getMaxDaysBeforeLatestDate())
	if err != nil {
		return err
	}

	var allStocksInDb []*models.Stock
	allStocksInDb, err = m.StockRepository.GetAll()
	if err != nil {
		return err
	}

	for _, stock := range allStocks {
		newStock := &models.Stock{
			Name:  stock.SecId,
			Price: stock.CurPrice,
		}

		err = m.createOrUpdateStock(newStock)
		if err != nil {
			return err
		}

		allStocksInDb = removeStockByNameFromSlice(allStocksInDb, newStock.Name)
	}

	if len(allStocksInDb) > 0 {
		err = m.removeOldStocksFromDb(allStocksInDb)
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

	var allIndexesInDb []*models.Index
	allIndexesInDb, err = m.IndexRepository.GetAll()
	if err != nil {
		return err
	}

	for _, index := range allIndexes {
		var indexStocks []*moex_models.IndexStocksData
		indexStocks, err = m.parseLatestIndexStocks(index, getMaxDaysBeforeLatestDate())
		if err != nil {
			return err
		}
		if indexStocks == nil {
			continue // Index contains bonds OR it's weekend (you can access index names but not its stocks)
		}

		newStocksFractionMap := make(map[*models.Stock]float64)
		newStocksFractionMap, err = m.createStocksFractionMapFromIndexStocks(indexStocks)
		if err != nil {
			return err
		}

		newIndex := &models.Index{
			Name:              index.IndexId,
			StocksFractionMap: newStocksFractionMap,
		}

		err = m.createOrUpdateIndex(newIndex)
		if err != nil {
			return err
		}

		allIndexesInDb = removeIndexByNameFromSlice(allIndexesInDb, newIndex.Name)
	}

	if len(allIndexesInDb) > 0 {
		err = m.removeOldIndexesFromDb(allIndexesInDb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) setApiClient() {
	m.moexApiClient = moex_api_client.NewMoexApiClient()
}

func getMaxDaysBeforeLatestDate() int {
	return 15
}

func removeElementFromSliceByIndex[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func getCurrentTime() time.Time {
	currTime := time.Now()
	return currTime
}

func formatTime(time time.Time) string {
	return time.Format("2006-01-02")
}
