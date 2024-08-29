package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/drawer/chart"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

type RouteHandler struct {
	portfolioRepository repository.PortfolioRepository
	assetRepository     repository.AssetRepository
	indexRepository     repository.IndexRepository
}

func NewRouteHandler(p repository.PortfolioRepository,
	s repository.AssetRepository, i repository.IndexRepository) *RouteHandler {
	return &RouteHandler{portfolioRepository: p, assetRepository: s, indexRepository: i}
}

func (r *RouteHandler) HandleStaticFiles(w http.ResponseWriter, request *http.Request) {
	ext := filepath.Ext(request.URL.Path)
	switch ext {
	case ".js", ".mjs":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "text/html")
	}

	http.ServeFile(w, request, "src/web/frontend"+request.URL.Path)
}
func (r *RouteHandler) HandleHome(w http.ResponseWriter, _ *http.Request) {

	files := []string{
		"./src/web/frontend/templates/home.page.tmpl",
		"./src/web/frontend/templates/navbar.tmpl",
		"./src/web/frontend/templates/base.page.tmpl",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	portfolios, err := r.portfolioRepository.GetAll()
	if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["portfolios"] = portfolios

	err = templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RouteHandler) HandleAddPortfolio(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		files := []string{
			"./src/web/frontend/templates/portfolio_creator.page.tmpl",
			"./src/web/frontend/templates/navbar.tmpl",
			"./src/web/frontend/templates/base.page.tmpl",
		}

		templates, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		portfolios, err := r.portfolioRepository.GetAll()
		if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		assets, err := r.assetRepository.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})
		data["portfolios"] = portfolios
		data["assets"] = assets

		err = templates.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

	case "POST":
		err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		portfolioName := request.FormValue("portfolioName")
		assets := request.MultipartForm.Value["assets[]"]

		_, err = r.portfolioRepository.GetByName(portfolioName)
		if !errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, "Some portfolio already has this name! Use another one", http.StatusConflict)
			return
		} else if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		type myAsset struct {
			Name     string `json:"name"`
			Quantity string `json:"Quantity"`
		}

		curAsset := myAsset{}

		assetsMap := make(map[*model.Asset]int)
		for _, asset := range assets {
			err = json.Unmarshal([]byte(asset), &curAsset)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			assetFromDB, err := r.assetRepository.GetByName(curAsset.Name)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			num, err := strconv.Atoi(curAsset.Quantity)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			assetsMap[assetFromDB] = num
		}

		portfolio := &model.Portfolio{
			Name:              portfolioName,
			AssetsQuantityMap: assetsMap,
		}

		err = r.portfolioRepository.Create(portfolio)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
	}
}

func handleNoPortfolioChosen(w http.ResponseWriter, _ *http.Request, data map[string]interface{}) {
	files := []string{
		"./src/web/frontend/templates/portfolio_chooser.page.tmpl",
		"./src/web/frontend/templates/navbar.tmpl",
		"./src/web/frontend/templates/base.page.tmpl",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func (r *RouteHandler) HandleManager(w http.ResponseWriter, request *http.Request) {
	portfolios, err := r.portfolioRepository.GetAll()
	if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie, err := request.Cookie("current_portfolio")
	if (err != nil) && (!errors.Is(err, http.ErrNoCookie)) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["portfolios"] = portfolios

	if cookie == nil || cookie.Value == "" || errors.Is(err, http.ErrNoCookie) {
		handleNoPortfolioChosen(w, request, data)
		return
	}

	portfolio, err := r.portfolioRepository.GetByName(cookie.Value)
	if errors.Is(err, repository.ErrPortfolioNotFound) {
		handleNoPortfolioChosen(w, request, data)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pieChart, err := chart.GetAssetPieChart(portfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type curAsset struct {
		Name     string
		Quantity int
		Price    float64
	}

	var assets []*curAsset
	for asset, quantity := range portfolio.AssetsQuantityMap {
		assets = append(assets, &curAsset{
			Name:     asset.Name,
			Quantity: quantity,
			Price:    asset.Price,
		})
	}

	data["PieChart"] = base64.StdEncoding.EncodeToString(pieChart)
	data["assets"] = assets

	files := []string{
		"./src/web/frontend/templates/manager.page.tmpl",
		"./src/web/frontend/templates/navbar.tmpl",
		"./src/web/frontend/templates/base.page.tmpl",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RouteHandler) HandleRemovePortfolio(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		portfolioName := request.FormValue("portfolioName")
		if portfolioName == "" {
			http.Error(w, "Error! No portfolioName in query!", http.StatusBadRequest)
			return
		}

		portfolio, err := r.portfolioRepository.GetByName(portfolioName)
		if portfolio == nil {
			http.Error(w, "Error! No such portfolio exists!", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err = r.portfolioRepository.Delete(portfolio); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// TODO move this function to somewhere?????
func countPriceOfPortfolio(portfolio *model.Portfolio) float64 {
	var price float64 = 0
	for asset, quantity := range portfolio.AssetsQuantityMap {
		price += asset.Price * float64(quantity)
	}

	return price
}

func (r *RouteHandler) HandleFollowingIndex(w http.ResponseWriter, request *http.Request) {
	portfolios, err := r.portfolioRepository.GetAll()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["portfolios"] = portfolios

	cookie, err := request.Cookie("current_portfolio")
	if (err != nil) && (!errors.Is(err, http.ErrNoCookie)) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cookie == nil || cookie.Value == "" || errors.Is(err, http.ErrNoCookie) {
		handleNoPortfolioChosen(w, request, data)
		return
	}

	_, err = r.portfolioRepository.GetByName(cookie.Value)
	if errors.Is(err, repository.ErrPortfolioNotFound) {
		handleNoPortfolioChosen(w, request, data)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indexes, err := r.indexRepository.GetAll()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data["indexes"] = indexes
	//data["budget"] = countPriceOfPortfolio(portfolio)

	files := []string{
		"./src/web/frontend/templates/following_index.page.tmpl",
		"./src/web/frontend/templates/navbar.tmpl",
		"./src/web/frontend/templates/base.page.tmpl",
	}

	templates, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = templates.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO move this function to somewhere?????
func tryFindAssetInAssetsMap[V any](mp *map[*model.Asset]V, name string) *model.Asset {
	for asset := range *mp {
		if asset.Name == name {
			return asset
		}
	}

	return nil
}

func (r *RouteHandler) HandleRenderFollowingIndexTable(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		indexForm := request.FormValue("index")
		//budget := request.FormValue("budget")
		portfolioName := request.FormValue("portfolio")

		// TODO check if indexForm, budget and portfolio name are valid

		portfolio, err := r.portfolioRepository.GetByName(portfolioName)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if portfolio == nil {
			http.Error(w, "No portfolio found", http.StatusBadRequest)
			return
		}

		index, err := r.indexRepository.GetByName(indexForm)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if index == nil {
			http.Error(w, "No index found", http.StatusBadRequest)
			return
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

		var curAsset []*tableAssets

		for asset, fraction := range index.AssetsFractionMap {
			tableAsset := &tableAssets{
				Name:              asset.Name,
				SuggestedFraction: fraction,
			}

			foundAsset := tryFindAssetInAssetsMap(&portfolio.AssetsQuantityMap, asset.Name)
			if foundAsset != nil {
				tableAsset.Quantity = portfolio.AssetsQuantityMap[foundAsset]
				tableAsset.Price = foundAsset.Price * float64(tableAsset.Quantity)
				tableAsset.CurrentFraction = 100 * float64(tableAsset.Quantity) * asset.Price / countPriceOfPortfolio(portfolio)
			} else {
				tableAsset.Quantity = 0
				tableAsset.Price = 0
				tableAsset.CurrentFraction = 0
			}

			tableAsset.DifferenceInFraction = tableAsset.CurrentFraction - tableAsset.SuggestedFraction

			curAsset = append(curAsset, tableAsset)
		}

		for asset := range portfolio.AssetsQuantityMap {
			foundInIndex := tryFindAssetInAssetsMap(&index.AssetsFractionMap, asset.Name)
			if foundInIndex != nil {
				continue
			}

			curFraction := 100 * float64(portfolio.AssetsQuantityMap[asset]) * asset.Price / countPriceOfPortfolio(portfolio)

			tableAsset := &tableAssets{
				Name:                 asset.Name,
				Quantity:             portfolio.AssetsQuantityMap[asset],
				Price:                asset.Price,
				CurrentFraction:      curFraction,
				SuggestedFraction:    0,
				DifferenceInFraction: curFraction,
			}

			curAsset = append(curAsset, tableAsset)
		}

		data := make(map[string]interface{})
		data["chart"] = "todo_chart_base64_here"
		data["assets"] = curAsset
		files := []string{
			"./src/web/frontend/templates/following_index.table.tmpl",
		}

		templates, err := template.ParseFiles(files...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = templates.Execute(w, data); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
