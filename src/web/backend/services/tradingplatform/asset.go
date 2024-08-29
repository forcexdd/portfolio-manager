package tradingplatform

import (
	"errors"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	moexmodels "github.com/forcexdd/portfoliomanager/src/web/backend/services/tradingplatform/moex/model"
	"time"
)

func (m *moexService) parseLatestAssets(maxDays int) ([]*moexmodels.AssetData, time.Time, error) {
	parseTime := getCurrentTime()

	allAssets, err := m.moexApiClient.GetAllAssets(formatTime(parseTime))
	if err != nil {
		return nil, time.Time{}, err
	}
	maxDays--

	for len(allAssets) == 0 && maxDays != 0 {
		parseTime = parseTime.AddDate(0, 0, -1)

		allAssets, err = m.moexApiClient.GetAllAssets(formatTime(parseTime))
		if err != nil {
			return nil, time.Time{}, err
		}

		maxDays--
	}

	if maxDays == 0 {
		m.log.Error("Exceeded limit on max days before latest date")
		return nil, time.Time{}, errors.New("exceeded limit on max days before latest date")
	}

	m.log.Info("Parsing date", "date", formatTime(parseTime))

	return allAssets, parseTime, nil
}

func (m *moexService) createOrUpdateAsset(asset *model.Asset) error {
	_, err := m.AssetRepository.GetByName(asset.Name)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) { // Can't find asset in database
			err = m.AssetRepository.Create(asset)
			if err != nil {
				return err
			}
			m.log.Info("Created asset", "name", asset.Name)
		}
		return err
	} else { // Found already existing asset in database
		err = m.AssetRepository.Update(asset)
		if err != nil {
			return err
		}
		m.log.Info("Updated asset", "name", asset.Name)
	}

	return nil
}

func removeAssetByNameFromSlice(assets []*model.Asset, name string) []*model.Asset {
	for index, asset := range assets {
		if asset.Name == name {
			return removeElementFromSliceByIndex[*model.Asset](assets, index)
		}
	}

	return assets
}

func (m *moexService) removeOldAssetsFromDB(assets []*model.Asset) error {
	for _, asset := range assets {
		err := m.AssetRepository.Delete(asset)
		if err != nil {
			return err
		}
		m.log.Info("Removed asset", "name", asset.Name)
	}

	return nil
}
