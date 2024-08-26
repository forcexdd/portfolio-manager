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

	assetRepository := repository.NewAssetRepository(db.GetDb())
	indexRepository := repository.NewIndexRepository(db.GetDb())

	tradingPlatformService := tradingplatform.NewTradingPlatformService(assetRepository, indexRepository)

	err = tradingPlatformService.ParseAllAssetsIntoDb()
	if err != nil {
		panic(err)
	}

	err = tradingPlatformService.ParseAllIndexesIntoDb()
	if err != nil {
		panic(err)
	}

	err = db.CloseConnection()
	if err != nil {
		panic(err)
	}

	log.Println("Time passed: ", time.Since(start))
}
