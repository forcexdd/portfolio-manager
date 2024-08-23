package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service"
	"log"
	"time"
)

func main() {
	start := time.Now()

	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		panic(err)
	}
	//db.DeleteStorage()

	stockRepository := repositories.NewStockRepository(db.GetDb())
	indexRepository := repositories.NewIndexRepository(db.GetDb())

	stockExchangeService := stock_exchange_service.NewStockExchangeService(stockRepository, indexRepository)

	err = stockExchangeService.ParseAllStocksIntoDb()
	if err != nil {
		panic(err)
	}

	err = stockExchangeService.ParseAllIndexesIntoDb()
	if err != nil {
		panic(err)
	}

	err = db.CloseConnection()
	if err != nil {
		panic(err)
	}

	log.Println("Time passed: ", time.Since(start))
}
