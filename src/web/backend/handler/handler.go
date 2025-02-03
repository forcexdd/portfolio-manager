package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/forcexdd/portfoliomanager/src/internal/logger"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/repository"
	"github.com/forcexdd/portfoliomanager/src/web/backend/database/storage"
	"github.com/forcexdd/portfoliomanager/src/web/backend/model"
	"github.com/forcexdd/portfoliomanager/src/web/backend/services/drawer/chart"
)

type RouteHandler struct {
	server              *http.Server
	db                  *storage.Storage
	portfolioRepository repository.PortfolioRepository
	assetRepository     repository.AssetRepository
	indexRepository     repository.IndexRepository
	log                 logger.Logger
}

func NewRouteHandler(server *http.Server, db *storage.Storage,
	p repository.PortfolioRepository, a repository.AssetRepository,
	i repository.IndexRepository, log logger.Logger) *RouteHandler {
	return &RouteHandler{
		server:              server,
		db:                  db,
		portfolioRepository: p,
		assetRepository:     a,
		indexRepository:     i,
		log:                 log,
	}
}

func (r *RouteHandler) HandleShutdown(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.db.CloseConnection(); err != nil {
		return
	}

	if err := r.server.Close(); err != nil {
		r.log.Error("Server failed to close", "error", err)
		return
	}
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

	http.ServeFile(w, request, frontendFolderPath+request.URL.Path)
}

func (r *RouteHandler) HandleHome(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		portfolios, err := r.portfolioRepository.GetAll()
		if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]any)
		data[allPortfoliosKey] = portfolios

		files := []string{
			homePageTemplatePath,
			navbarTemplatePath,
			basePageTemplatePath,
		}

		if err = renderTemplates(w, files, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (r *RouteHandler) HandleGetAddPortfolio(w http.ResponseWriter) {
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

	data := make(map[string]any)
	data[allPortfoliosKey] = portfolios
	data[allAssetsKey] = assets

	files := []string{
		portfolioCreatorPageTemplatePath,
		navbarTemplatePath,
		basePageTemplatePath,
	}

	if err = renderTemplates(w, files, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *RouteHandler) HandlePostAddPortfolio(w http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	portfolioName := request.FormValue(portfolioFormKey)
	assets := request.MultipartForm.Value[allAssetsFormKey]

	_, err = r.portfolioRepository.GetByName(portfolioName)
	if !errors.Is(err, repository.ErrPortfolioNotFound) {
		http.Error(w, "Some portfolio already has this name! Use another one", http.StatusConflict)
		return
	} else if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assetsMap, err := convertFormAssetsToModelAssetsMap(assets, r.assetRepository)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	portfolio := &model.Portfolio{
		Name:              portfolioName,
		AssetsQuantityMap: *assetsMap,
	}

	err = r.portfolioRepository.Create(portfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *RouteHandler) HandleAddPortfolio(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		r.HandleGetAddPortfolio(w)

	case http.MethodPost:
		r.HandlePostAddPortfolio(w, request)
	}
}

func handleNoPortfolioChosen(w http.ResponseWriter, _ *http.Request, data map[string]any) {
	files := []string{
		portfolioChooserPageTemplatePath,
		navbarTemplatePath,
		basePageTemplatePath,
	}

	if err := renderTemplates(w, files, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *RouteHandler) HandleManager(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		portfolios, err := r.portfolioRepository.GetAll()
		if err != nil && !errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]any)
		data[allPortfoliosKey] = portfolios

		cookie, err := request.Cookie(portfolioCookieName)
		if (err != nil) && (!errors.Is(err, http.ErrNoCookie)) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		assets := convertModelAssetsToHTMLAssets(&portfolio.AssetsQuantityMap)

		data[pieChartKey] = base64.StdEncoding.EncodeToString(pieChart)
		data[allAssetsKey] = *assets

		files := []string{
			managerPageTemplatePath,
			navbarTemplatePath,
			basePageTemplatePath,
		}

		if err = renderTemplates(w, files, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (r *RouteHandler) HandleRemovePortfolio(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		portfolioName := request.FormValue(portfolioFormKey)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = r.portfolioRepository.Delete(portfolio); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (r *RouteHandler) HandleFollowingIndex(w http.ResponseWriter, request *http.Request) {
	portfolios, err := r.portfolioRepository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]any)
	data[allPortfoliosKey] = portfolios

	cookie, err := request.Cookie(portfolioCookieName)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data[allIndexesKey] = indexes
	//data["budget"] = countPriceOfPortfolio(portfolio)

	files := []string{
		followingIndexPageTemplatePath,
		navbarTemplatePath,
		basePageTemplatePath,
	}

	if err = renderTemplates(w, files, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *RouteHandler) HandleRenderFollowingIndexTable(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		err := request.ParseMultipartForm(-1) // we can use -1 for a no memory limit ???
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		indexForm := request.FormValue(indexFormKey)
		//budget := request.FormValue("budget")
		portfolioName := request.FormValue(portfolioFormKey)

		// TODO check if indexForm, budget and portfolio name are valid

		portfolio, err := r.portfolioRepository.GetByName(portfolioName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if portfolio == nil {
			http.Error(w, "No portfolio found", http.StatusBadRequest)
			return
		}

		index, err := r.indexRepository.GetByName(indexForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if index == nil {
			http.Error(w, "No index found", http.StatusBadRequest)
			return
		}

		portfolioPrice := countPriceOfPortfolio(portfolio)
		curAssets := getHTMLTableAssetsFromIndex(&index.AssetsFractionMap, &portfolio.AssetsQuantityMap, portfolioPrice)

		// Get assets from a user portfolio that are not in index
		usersUnusedInIndexAssets := getUsersNotInIndexAssetsAsHTMLTableAssets(&index.AssetsFractionMap, &portfolio.AssetsQuantityMap, portfolioPrice)

		*curAssets = append(*curAssets, *usersUnusedInIndexAssets...)
		sort.Sort(byName(*curAssets))

		data := make(map[string]any)
		//data["chart"] = "todo_chart_base64_here" TODO
		data[allAssetsKey] = *curAssets

		files := []string{
			followingIndexTableTemplatePath,
		}

		if err = renderTemplates(w, files, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (r *RouteHandler) HandleAddAsset(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := request.Cookie(portfolioCookieName)
	if err != nil || cookie.Value == "" {
		http.Error(w, "No portfolio selected", http.StatusBadRequest)
		return
	}

	portfolio, err := r.portfolioRepository.GetByName(cookie.Value)
	if err != nil {
		if errors.Is(err, repository.ErrPortfolioNotFound) {
			http.Error(w, "Portfolio not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	assetName := request.FormValue("assetName")
	if assetName == "" {
		http.Error(w, "Asset name required", http.StatusBadRequest)
		return
	}

	for portolioAsset := range portfolio.AssetsQuantityMap {
		if portolioAsset.Name == assetName {
			portfolio.AssetsQuantityMap[portolioAsset]++
		}
	}

	if err := r.portfolioRepository.Update(portfolio); err != nil {
		http.Error(w, "Failed to update portfolio", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, request, request.Referer(), http.StatusSeeOther)
}

func (r *RouteHandler) HandleRemoveAsset(w http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := request.Cookie(portfolioCookieName)
	if err != nil || cookie.Value == "" {
		http.Error(w, "No portfolio selected", http.StatusBadRequest)
		return
	}

	portfolio, err := r.portfolioRepository.GetByName(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := request.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	assetName := request.FormValue("assetName")
	if assetName == "" {
		return
	}

	for portolioAsset, currentQuantity := range portfolio.AssetsQuantityMap {
		if portolioAsset.Name == assetName {
			if currentQuantity > 1 {
				portfolio.AssetsQuantityMap[portolioAsset]--
			} else {
				delete(portfolio.AssetsQuantityMap, portolioAsset)
			}
		}
	}

	if err := r.portfolioRepository.Update(portfolio); err != nil {
		return
	}

	http.Redirect(w, request, request.Referer(), http.StatusSeeOther)
}
