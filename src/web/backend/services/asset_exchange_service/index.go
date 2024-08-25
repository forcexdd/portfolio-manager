package asset_exchange_service

import (
	"github.com/forcexdd/portfolio_manager/src/web/backend/models"
	"github.com/forcexdd/portfolio_manager/src/web/backend/services/asset_exchange_service/moex/moex_models"
)

func (m *MoexService) parseLatestIndexAssets(index *moex_models.IndexData, maxDays int) ([]*moex_models.IndexAssetsData, error) {
	parseTime := getCurrentTime()

	allIndexAssets, err := m.moexApiClient.GetAllIndexAssets(formatTime(parseTime), index)
	if err != nil {
		return nil, err
	}
	maxDays--

	for len(allIndexAssets) == 0 && maxDays != 0 {
		parseTime = parseTime.AddDate(0, 0, -1)

		allIndexAssets, err = m.moexApiClient.GetAllIndexAssets(formatTime(parseTime), index)
		if err != nil {
			return nil, err
		}

		maxDays--
	}

	return allIndexAssets, nil
}

func (m *MoexService) createOrUpdateIndex(index *models.Index) error {
	dbIndex, err := m.IndexRepository.GetByName(index.Name)
	if err != nil {
		return err

	} else if dbIndex == nil { // Can't find index in database
		err = m.IndexRepository.Create(index)
		if err != nil {
			return err
		}
	} else { // Found already existing index in database
		err = m.IndexRepository.Update(index)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MoexService) createAssetsFractionMapFromIndexAssets(indexAssets []*moex_models.IndexAssetsData) (map[*models.Asset]float64, error) {
	newAssetsFractionMap := make(map[*models.Asset]float64)
	for _, indexAsset := range indexAssets {
		asset, err := m.AssetRepository.GetByName(indexAsset.SecIds)
		if err != nil {
			return nil, err
		}

		newAssetsFractionMap[asset] = indexAsset.Weight
	}

	return newAssetsFractionMap, nil
}

func removeIndexByNameFromSlice(indexes []*models.Index, name string) []*models.Index {
	for i, index := range indexes {
		if index.Name == name {
			return removeElementFromSliceByIndex[*models.Index](indexes, i)
		}
	}

	return indexes
}

func (m *MoexService) removeOldIndexesFromDb(indexes []*models.Index) error {
	for _, index := range indexes {
		err := m.IndexRepository.Delete(index)
		if err != nil {
			return err
		}
	}

	return nil
}
