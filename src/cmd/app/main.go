package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/handlers"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/static/", handlers.HandleStaticFiles)
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/following_index", handlers.HandleHome)
	http.HandleFunc("/manager", handlers.HandleManager)
	http.HandleFunc("/add_portfolio", handlers.HandleAddPortfolio)

	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		log.Fatal(err)
	}

	stockRepository := repositories.NewStockRepository(db.GetDb())
	portfolioRepository := repositories.NewPortfolioRepository(db.GetDb())

	stock := &models.Stock{Name: "AFLT", Price: 123.4567}
	stockRepository.Delete(stock)

	//newStock, err := stockRepository.GetByName("AFLT")
	newMap := map[*models.Stock]int{
		stock: 3,
	}

	portfolioRepository.Delete(&models.Portfolio{Name: "Aboba portfel", StocksQuantityMap: newMap})

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	
	db.CloseConnection()
}
