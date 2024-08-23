package moex_service

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/moex_api_client"
)

type StockExchangeService interface {
	ParseAllStocksIntoDb()
	ParseAllIndexesIntoDb()
}

type MoexService struct {
	StockRepository repositories.StockRepository
	IndexRepository repositories.IndexRepository
	MoexApiClient   moex_api_client.MoexApiClient
}

func NewStockExchangeService(stockRepository repositories.StockRepository, indexRepository repositories.IndexRepository, moexApiClient moex_api_client.MoexApiClient) StockExchangeService {
	return &MoexService{
		StockRepository: stockRepository,
		IndexRepository: indexRepository,
		MoexApiClient:   moexApiClient,
	}
}

func (m *MoexService) ParseAllStocksIntoDb() {}
