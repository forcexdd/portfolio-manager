package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"log"
)

func main() {
	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		log.Fatal(err)
	}

	//stockRepository := repositories.NewStockRepository(db.GetDb())
	//portfolioRepository := repositories.NewPortfolioRepository(db.GetDb())
	//
	//stockAFLT := &models.Stock{Name: "AFLT", Price: 123.4567}
	//stockGAZP := &models.Stock{Name: "GAZP", Price: 456}
	//stockALRS := &models.Stock{Name: "ALRS", Price: 45.12}
	//newMap := map[*models.Stock]int{
	//	stockAFLT: 1,
	//	stockGAZP: 34,
	//	stockALRS: 9,
	//}
	//
	//newPortfolio := &models.Portfolio{Name: "Aboba", StocksQuantityMap: newMap}
	//
	//stockRepository.Create(stockAFLT)
	//stockRepository.Create(stockGAZP)
	//stockRepository.Create(stockALRS)
	//
	//portfolioRepository.Create(newPortfolio)

	db.CloseConnection()
}
