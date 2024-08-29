package main

import (
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform"
	"log"
	"time"
)

func main() {
	start := time.Now()

	const connString = "postgresql://postgres:postgres@localhost:5432/portfolio_manager?sslmode=disable"

	db, err := storage.NewStorage(connString)
	if err != nil {
		panic(err)
	}
	//db.DeleteStorage()

	assetRepository := repository.NewAssetRepository(db.GetDB())
	indexRepository := repository.NewIndexRepository(db.GetDB())

	tradingPlatformService := tradingplatform.NewTradingPlatformService(assetRepository, indexRepository)

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

	log.Println("Time passed: ", time.Since(start))
}
