package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/database/repositories"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services/chart_drawer"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
)

type RouteHandler struct {
	portfolioRepository repositories.PortfolioRepository
	stocksRepository    repositories.StockRepository
	indexRepository     repositories.IndexRepository
}

func NewRouteHandler(p repositories.PortfolioRepository,
	s repositories.StockRepository, i repositories.IndexRepository) *RouteHandler {
	return &RouteHandler{portfolioRepository: p, stocksRepository: s, indexRepository: i}
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
	if err != nil {
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stocks, err := r.stocksRepository.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := make(map[string]interface{})
		data["portfolios"] = portfolios
		data["stocks"] = stocks

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
		stocks := request.MultipartForm.Value["stocks[]"]

		foundPortfolio, err := r.portfolioRepository.GetByName(portfolioName)
		if foundPortfolio != nil {
			http.Error(w, "Some portfolio already has this name! Use another one", http.StatusConflict)
			return
		}

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		type myStock struct {
			Name     string `json:"name"`
			Quantity string `json:"Quantity"`
		}

		curStock := myStock{}

		stocksMap := make(map[*models.Stock]int)
		for _, stock := range stocks {
			err = json.Unmarshal([]byte(stock), &curStock)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			stockFromDb, err := r.stocksRepository.GetByName(curStock.Name)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			num, err := strconv.Atoi(curStock.Quantity)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			stocksMap[stockFromDb] = num
		}

		portfolio := &models.Portfolio{
			Name:              portfolioName,
			StocksQuantityMap: stocksMap,
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
	if err != nil {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if portfolio == nil {
		handleNoPortfolioChosen(w, request, data)
		return
	}

	pieChart, err := chart_drawer.GetStockPieChart(portfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type curStock struct {
		Name     string
		Quantity int
		Price    float64
	}

	var stocks []*curStock
	for stock, quantity := range portfolio.StocksQuantityMap {
		stocks = append(stocks, &curStock{
			Name:     stock.Name,
			Quantity: quantity,
			Price:    stock.Price,
		})
	}

	data["PieChart"] = base64.StdEncoding.EncodeToString(pieChart)
	data["stocks"] = stocks

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
func countPriceOfPortfolio(portfolio *models.Portfolio) float64 {
	var price float64 = 0
	for stock, quantity := range portfolio.StocksQuantityMap {
		price += stock.Price * float64(quantity)
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

	portfolio, err := r.portfolioRepository.GetByName(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if portfolio == nil {
		handleNoPortfolioChosen(w, request, data)
		return
	}

	indexes, err := r.indexRepository.GetAll()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data["indexes"] = indexes
	data["budget"] = countPriceOfPortfolio(portfolio)

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

		type tableStocks struct {
			Name                 string
			Quantity             int
			Price                float64
			CurrentFraction      float64
			SuggestedFraction    float64
			DifferenceInFraction float64
			//Advice            string
		}

		var curStocks []*tableStocks

		// TODO move this function to somewhere?????
		tryFindStockInPortfolio := func(portfolio *models.Portfolio, name string) *models.Stock {
			for stock := range portfolio.StocksQuantityMap {
				if stock.Name == name {
					return stock
				}
			}

			return nil
		}

		for stock, fraction := range index.StocksFractionMap {
			tableStock := &tableStocks{
				Name:              stock.Name,
				SuggestedFraction: fraction,
			}

			foundStock := tryFindStockInPortfolio(portfolio, stock.Name)
			if foundStock != nil {
				tableStock.Quantity = portfolio.StocksQuantityMap[foundStock]
				tableStock.Price = foundStock.Price * float64(tableStock.Quantity)
				tableStock.CurrentFraction = 100 * float64(tableStock.Quantity) * stock.Price / countPriceOfPortfolio(portfolio)
			} else {
				tableStock.Quantity = 0
				tableStock.Price = 0
				tableStock.CurrentFraction = 0
			}

			tableStock.DifferenceInFraction = tableStock.CurrentFraction - tableStock.SuggestedFraction

			curStocks = append(curStocks, tableStock)
		}

		data := make(map[string]interface{})
		data["chart"] = "todo_chart_base64_here"
		data["stocks"] = curStocks
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
