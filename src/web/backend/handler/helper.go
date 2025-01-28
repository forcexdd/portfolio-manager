package handler

import (
	"encoding/json"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	"net/http"
	"strconv"
	"text/template"
)

// TODO rename package and move some functions somewhere

func renderTemplates(w http.ResponseWriter, files []string, data map[string]any) error {
	templates, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}

	err = templates.Execute(w, data)

	return err
}

func countPriceOfPortfolio(portfolio *model.Portfolio) float64 {
	var price float64 = 0
	for asset, quantity := range portfolio.AssetsQuantityMap {
		price += asset.Price * float64(quantity)
	}

	return price
}

func tryFindAssetInAssetsMap[V any](mp *map[*model.Asset]V, name string) *model.Asset {
	for asset := range *mp {
		if asset.Name == name {
			return asset
		}
	}

	return nil
}

func convertFormAssetsToModelAssetsMap(formAssets []string, assetRepository repository.AssetRepository) (*map[*model.Asset]int, error) {
	// Shit func, remove assetRepository but how to make it work then?
	type formAsset struct {
		Name     string `json:"name"`
		Quantity string `json:"Quantity"`
	}

	curAsset := formAsset{}

	assetsMap := make(map[*model.Asset]int)
	for _, asset := range formAssets {
		if err := json.Unmarshal([]byte(asset), &curAsset); err != nil {
			return nil, err
		}

		assetFromDB, err := assetRepository.GetByName(curAsset.Name)
		if err != nil {
			return nil, err
		}

		num, err := strconv.Atoi(curAsset.Quantity)
		if err != nil {
			return nil, err
		}
		assetsMap[assetFromDB] = num
	}

	return &assetsMap, nil
}

type HTMLAsset struct {
	Name     string
	Quantity int
	Price    float64
}

func convertModelAssetsToHTMLAssets(assetsQuantityMap *map[*model.Asset]int) *[]*HTMLAsset {
	var assets []*HTMLAsset

	for asset, quantity := range *assetsQuantityMap {
		assets = append(assets, &HTMLAsset{
			Name:     asset.Name,
			Quantity: quantity,
			Price:    asset.Price,
		})
	}

	return &assets
}

type tableAssets struct {
	Name                 string
	Quantity             int
	Price                float64
	CurrentFraction      float64
	SuggestedFraction    float64
	DifferenceInFraction float64
	//Advice            string
}

func getHTMLTableAssetsFromIndex(indexAssetsMap *map[*model.Asset]float64,
	portfolioAssetsMap *map[*model.Asset]int, portfolioPrice float64) *[]*tableAssets {
	var curAsset []*tableAssets

	for asset, fraction := range *indexAssetsMap {
		tableAsset := &tableAssets{
			Name:              asset.Name,
			SuggestedFraction: fraction,
		}

		foundAsset := tryFindAssetInAssetsMap(portfolioAssetsMap, asset.Name)
		if foundAsset != nil {
			tableAsset.Quantity = (*portfolioAssetsMap)[foundAsset]
			tableAsset.Price = foundAsset.Price * float64(tableAsset.Quantity)
			tableAsset.CurrentFraction = 100 * float64(tableAsset.Quantity) * asset.Price / portfolioPrice
		} else {
			tableAsset.Quantity = 0
			tableAsset.Price = 0
			tableAsset.CurrentFraction = 0
		}

		tableAsset.DifferenceInFraction = tableAsset.CurrentFraction - tableAsset.SuggestedFraction

		curAsset = append(curAsset, tableAsset)
	}

	return &curAsset
}

func getUsersNotInIndexAssetsAsHTMLTableAssets(indexAssetsMap *map[*model.Asset]float64,
	portfolioAssetsMap *map[*model.Asset]int, portfolioPrice float64) *[]*tableAssets {
	var curAssets []*tableAssets

	for asset := range *portfolioAssetsMap {
		foundInIndex := tryFindAssetInAssetsMap(indexAssetsMap, asset.Name)
		if foundInIndex != nil {
			continue
		}

		curFraction := 100 * float64((*portfolioAssetsMap)[asset]) * asset.Price / portfolioPrice

		tableAsset := &tableAssets{
			Name:                 asset.Name,
			Quantity:             (*portfolioAssetsMap)[asset],
			Price:                asset.Price,
			CurrentFraction:      curFraction,
			SuggestedFraction:    0,
			DifferenceInFraction: curFraction,
		}

		curAssets = append(curAssets, tableAsset)
	}

	return &curAssets
}

type byName []*tableAssets

func (a byName) Len() int {
	return len(a)
}

func (a byName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}
