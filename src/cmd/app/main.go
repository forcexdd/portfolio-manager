package main

import (
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/handler"
	"log"
	"net/http"
	"time"
)

func main() {
	start := time.Now()

	const connString = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		panic(err)
	}
	//err = db.DeleteStorage()
	//if err != nil {
	//	panic(err)
	//}

	assetRepository := repository.NewAssetRepository(db.GetDb())
	indexRepository := repository.NewIndexRepository(db.GetDb())
	portfolioRepository := repository.NewPortfolioRepository(db.GetDb())

	routeHandler := handler.NewRouteHandler(portfolioRepository, assetRepository, indexRepository)

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

	log.Println("Time passed: ", time.Since(start))
}
