package main

import (
	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform"
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
	//db.DeleteStorage()

	assetRepository := repository.NewAssetRepository(db.GetDB(), log)
	indexRepository := repository.NewIndexRepository(db.GetDB(), log)

	tradingPlatformService := tradingplatform.NewTradingPlatformService(assetRepository, indexRepository, log)

	err = tradingPlatformService.ParseAllAssetsIntoDB()
	if err != nil {
		panic(err)
	}

	err = tradingPlatformService.ParseAllIndexesIntoDB()
	if err != nil {
		panic(err)
	}

	err = db.CloseConnection()
	if err != nil {
		panic(err)
	}

	err = log.Close()
	if err != nil {
		panic(err)
	}

	log.Info("App closed", "time passed", time.Since(start))
}
