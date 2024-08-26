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
	indexRepository := repositories.NewIndexRepository(db.GetDb())
	portfolioRepository := repositories.NewPortfolioRepository(db.GetDb())

	routeHandler := handlers.NewRouteHandler(portfolioRepository, stockRepository, indexRepository)

	http.HandleFunc("/static/", routeHandler.HandleStaticFiles)
	http.HandleFunc("/", routeHandler.HandleHome)
	http.HandleFunc("/following_index", routeHandler.HandleFollowingIndex)
	http.HandleFunc("/manager", routeHandler.HandleManager)
	http.HandleFunc("/add_portfolio", routeHandler.HandleAddPortfolio)
	http.HandleFunc("/remove_portfolio", routeHandler.HandleRemovePortfolio)
	http.HandleFunc("/render_following_index_table", routeHandler.HandleRenderFollowingIndexTable)
	if err = http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}

	if err = db.CloseConnection(); err != nil {
		panic(err)
	}
}
