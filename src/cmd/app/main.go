package main

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/storage"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/handlers"
	"net/http"
)

func main() {
	const connString = "postgresql://postgres:postgres@localhost:5432/StockPortfolioManager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		panic(err)
	}
	//err = db.DeleteStorage()
	//if err != nil {
	//	panic(err)
	//}

	stockRepository := repositories.NewStockRepository(db.GetDb())
	portfolioRepository := repositories.NewPortfolioRepository(db.GetDb())

	routeHandler := handlers.NewRouteHandler(portfolioRepository, stockRepository)

	http.HandleFunc("/static/", routeHandler.HandleStaticFiles)
	http.HandleFunc("/", routeHandler.HandleHome)
	http.HandleFunc("/following_index", routeHandler.HandleHome)
	http.HandleFunc("/manager", routeHandler.HandleManager)
	http.HandleFunc("/add_portfolio", routeHandler.HandleAddPortfolio)
	if err = http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}

	if err = db.CloseConnection(); err != nil {
		panic(err)
	}
}
