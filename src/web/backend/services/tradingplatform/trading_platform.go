package tradingplatform

import (
	"errors"
	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/client"
	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"time"
)

type AssetExchangeService interface {
	ParseAllAssetsIntoDB() error
	ParseAllIndexesIntoDB() error
}

type moexService struct {
	AssetRepository repository.AssetRepository
	IndexRepository repository.IndexRepository
	moexApiClient   *client.MoexApiClient
	time            time.Time
	log             logger.Logger
}

func NewTradingPlatformService(assetRepository repository.AssetRepository, indexRepository repository.IndexRepository, log logger.Logger) AssetExchangeService {
	newService := &moexService{
		AssetRepository: assetRepository,
		IndexRepository: indexRepository,
		time:            time.Time{},
		log:             log,
	}
	newService.setApiClient()

	return newService
}

func (m *moexService) ParseAllAssetsIntoDB() error {
	var allAssets []*moexmodels.AssetData
	var err error
	allAssets, m.time, err = m.parseLatestAssets(minDaysBeforeLatestDate, maxDaysBeforeLatestDate)
	if err != nil {
		m.log.Error("Failed parsing latest assets", "error", err)
		return err
	}

	var allAssetsInDB []*model.Asset
	allAssetsInDB, err = m.AssetRepository.GetAll()
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) { // If we don't find any assets in DB it's alright
			err = nil
		} else {
			m.log.Error("Failed getting all assets from DB", "error", err)
			return err
		}
	}

	for _, asset := range allAssets {
		newAsset := &model.Asset{
			Name:  asset.SecID,
			Price: asset.CurPrice,
		}

		err = m.createOrUpdateAsset(newAsset)
		if err != nil {
			return err
		}

		allAssetsInDB = removeAssetByNameFromSlice(allAssetsInDB, newAsset.Name)
	}

	if len(allAssetsInDB) > 0 {
		err = m.removeOldAssetsFromDB(allAssetsInDB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *moexService) ParseAllIndexesIntoDB() error {
	if m.time.IsZero() {
		m.log.Warn("All assets have not been parsed yet")
		return errors.New("assets should be parsed first")
	}

	allIndexes, err := m.moexApiClient.GetAllIndexes()
	if err != nil {
		m.log.Error("Failed getting all indexes", "error", err)
		return err
	}

	var allIndexesInDB []*model.Index
	allIndexesInDB, err = m.IndexRepository.GetAll()
	if err != nil {
		if errors.Is(err, repository.ErrIndexNotFound) { // If we don't find any indexes in DB it's alright
			err = nil
		} else {
			return err
		}
	}

	for _, index := range allIndexes {
		var indexAssets []*moexmodels.IndexAssetsData
		indexAssets, err = m.parseIndexAssets(index, m.time)
		if err != nil {
			m.log.Error("Failed parsing assets associated with index", "name", index.IndexID, "error", err)
			return err
		}
		if len(indexAssets) != 0 { // No assets in index (something wrong with index API response so we skip this particular index)
			m.log.Warn("No assets in index", "name", index.IndexID)
			continue
		}

		newAssetsFractionMap := make(map[*model.Asset]float64)
		newAssetsFractionMap, err = m.createAssetsFractionMapFromIndexAssets(indexAssets)
		if err != nil {
			if errors.Is(err, repository.ErrAssetNotFound) { // There is no such asset in database (something wrong with index API response so we skip this particular index)
				m.log.Warn("No asset from index in DB", "name", index.IndexID)
				continue
			}
			return err
		}

		newIndex := &model.Index{
			Name:              index.IndexID,
			AssetsFractionMap: newAssetsFractionMap,
		}

		err = m.createOrUpdateIndex(newIndex)
		if err != nil {
			return err
		}

		allIndexesInDB = removeIndexByNameFromSlice(allIndexesInDB, newIndex.Name)

	}

	if len(allIndexesInDB) > 0 {
		err = m.removeOldIndexesFromDB(allIndexesInDB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *moexService) setApiClient() {
	m.moexApiClient = client.NewMoexApiClient()
}

func removeElementFromSliceByIndex[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func getCurrentTime() time.Time {
	currTime := time.Now()
	return currTime
}

func formatTime(time time.Time) string {
	return time.Format("2006-01-02")
}
