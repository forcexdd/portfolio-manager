package main

import (
	"net/http"

	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/handler"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform"
)

const (
	connString = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"
	logPath    = "portfolio_manager.log"
	url        = "localhost:8080"
)

func main() {
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

	server := &http.Server{
		Addr: url,
	}

	assetRepository := repository.NewAssetRepository(db.GetDB(), log)
	indexRepository := repository.NewIndexRepository(db.GetDB(), log)
	portfolioRepository := repository.NewPortfolioRepository(db.GetDB(), log)

	tradingPlatformService := tradingplatform.NewTradingPlatformService(
		assetRepository,
		indexRepository,
		log,
	)

	parsingErr := make(chan error)
	go func() {
		err := tradingPlatformService.ParseAllAssetsIntoDB()
		if err != nil {
			parsingErr <- err
			return
		}

		err = tradingPlatformService.ParseAllIndexesIntoDB()
		if err != nil {
			parsingErr <- err
			return
		}
	}()

	routeHandler := handler.NewRouteHandler(
		server,
		db,
		portfolioRepository,
		assetRepository,
		indexRepository,
		log,
	)

	http.HandleFunc("/shutdown", routeHandler.HandleShutdown)
	http.HandleFunc("/static/", routeHandler.HandleStaticFiles)
	http.HandleFunc("/", routeHandler.HandleHome)
	http.HandleFunc("/following_index", routeHandler.HandleFollowingIndex)
	http.HandleFunc("/manager", routeHandler.HandleManager)
	http.HandleFunc("/add_portfolio", routeHandler.HandleAddPortfolio)
	http.HandleFunc("/remove_portfolio", routeHandler.HandleRemovePortfolio)
	http.HandleFunc("/render_following_index_table", routeHandler.HandleRenderFollowingIndexTable)

	serverErr := make(chan error)
	go func() {
		log.Info("Server starting", "url", url)
		if err := server.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		switch {
		case err == http.ErrServerClosed:
			log.Info("Server closed")
		default:
			log.Error("Server failed to start", "error", err)
		}
	case err := <-parsingErr:
		log.Error("Parsing failed", "error", err)
	}

	log.Close()
}
