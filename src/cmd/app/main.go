package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/handlers"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/stock_exchange_service"
	"log"
	"net/http"
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
	portfolioRepository := repositories.NewPortfolioRepository(db.GetDb())
	indexRepository := repositories.NewIndexRepository(db.GetDb())

	stockExchangeService := stock_exchange_service.NewStockExchangeService(stockRepository, indexRepository)

	err = stockExchangeService.ParseAllStocksIntoDb()
	if err != nil {
		panic(err)
	}

	//stocksToCreate := []*models.Stock{
	//	{Name: "AFLT", Price: 123.4567},
	//	{Name: "ROSB", Price: 119.2834},
	//	{Name: "NKHP", Price: 874.0},
	//	{Name: "KLSB", Price: 27.3},
	//	{Name: "UNKL", Price: 5540.0},
	//}
	//
	//for _, stock := range stocksToCreate {
	//	stockRepository.Create(stock)
	//}
	//
	//mp := make(map[*models.Stock]int)
	//mp[stocksToCreate[0]] = 3
	//mp[stocksToCreate[2]] = 15
	//mp[stocksToCreate[3]] = 8
	//
	//portfolio := &models.Portfolio{
	//	Name:              "First",
	//	StocksQuantityMap: mp,
	//}
	//
	//portfolioRepository.Create(portfolio)

	routeHandler := handlers.NewRouteHandler(portfolioRepository, stockRepository)

	http.HandleFunc("/static/", routeHandler.HandleStaticFiles)
	http.HandleFunc("/", routeHandler.HandleHome)
	http.HandleFunc("/following_index", routeHandler.HandleHome)
	http.HandleFunc("/manager", routeHandler.HandleManager)
	http.HandleFunc("/add_portfolio", routeHandler.HandleAddPortfolio)

	//stock := &models.Stock{Name: "AFLT", Price: 123.4567}
	//stockRepository.Delete(stock)
	//
	////newStock, err := stockRepository.GetByName("AFLT")
	//newMap := map[*models.Stock]int{
	//	stock: 3,
	//}

	//portfolioRepository.Delete(&models.Portfolio{Name: "Aboba portfel", StocksQuantityMap: newMap})

	err = stockExchangeService.ParseAllIndexesIntoDb()
	if err != nil {
		panic(err)
	}
	log.Println("Time passed: ", time.Since(start))
	
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	err = db.CloseConnection()
	if err != nil {
		panic(err)
	}

}
