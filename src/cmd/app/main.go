package main

import (
	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/handler"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform"
	"net/http"
	"time"
)

const (
	connString = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"
	logPath    = "portfolio_manager.log"
)

func main() {
	start := time.Now()

	log, err := logger.NewLogger(logPath)
	if err != nil {
		panic(err)
	}

	var db *storage.Storage
	db, err = storage.NewStorage(connString, log)
	if err != nil {
		panic(err)
	}
	//err = db.DeleteStorage()
	//if err != nil {
	//	panic(err)

	assetRepository := repository.NewAssetRepository(db.GetDB(), log)
	indexRepository := repository.NewIndexRepository(db.GetDB(), log)
	portfolioRepository := repository.NewPortfolioRepository(db.GetDB(), log)

	tradingPlatformService := tradingplatform.NewTradingPlatformService(assetRepository, indexRepository, log)

	err = tradingPlatformService.ParseAllAssetsIntoDB()
	if err != nil {
		panic(err)
	}

	err = tradingPlatformService.ParseAllIndexesIntoDB()
	if err != nil {
		panic(err)
	}

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

	err = log.Close()
	if err != nil {
		panic(err)
	}

	log.Info("App closed", "time passed", time.Since(start))
}
