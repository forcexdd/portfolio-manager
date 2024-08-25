package asset_exchange_service

import (
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_models"
)

func (m *MoexService) parseLatestAssets(maxDays int) ([]*moex_models.AssetData, error) {
	parseTime := getCurrentTime()

	allAssets, err := m.moexApiClient.GetAllAssets(formatTime(parseTime))
	if err != nil {
		return nil, err
	}
	maxDays--

	for len(allAssets) == 0 && maxDays != 0 {
		parseTime = parseTime.AddDate(0, 0, -1)

		allAssets, err = m.moexApiClient.GetAllAssets(formatTime(parseTime))
		if err != nil {
			return nil, err
		}

		maxDays--
	}

	return allAssets, nil
}

func (m *MoexService) createOrUpdateAsset(asset *models.Asset) error {
	dbAsset, err := m.AssetRepository.GetByName(asset.Name)
	if err != nil {
		return err

	} else if dbAsset == nil { // Can't find asset in database
		err = m.AssetRepository.Create(asset)
		if err != nil {
			return err
		}
	} else { // Found already existing asset in database
		err = m.AssetRepository.Update(asset)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeAssetByNameFromSlice(assets []*models.Asset, name string) []*models.Asset {
	for index, asset := range assets {
		if asset.Name == name {
			return removeElementFromSliceByIndex[*models.Asset](assets, index)
		}
	}

	return assets
}

func (m *MoexService) removeOldAssetsFromDb(assets []*models.Asset) error {
	for _, asset := range assets {
		err := m.AssetRepository.Delete(asset)
		if err != nil {
			return err
		}
	}

	return nil
}
