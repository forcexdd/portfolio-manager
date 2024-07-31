package handlers

import (
	"encoding/base64"
	"errors"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/services"
	"net/http"
	"path/filepath"
	"text/template"
)

func HandleStaticFiles(w http.ResponseWriter, r *http.Request) {
	ext := filepath.Ext(r.URL.Path)
	switch ext {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "text/html")
	}

	http.ServeFile(w, r, "src/web/frontend"+r.URL.Path)
}
func HandleHome(w http.ResponseWriter, r *http.Request) {

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

	portfolios := []models.Portfolio{
		{"testFirst"},
		{"testSecond"},
	}

	data := make(map[string]interface{})
	data["portfolios"] = portfolios
	err = templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getData() map[string]interface{} {
	portfolios := []models.Portfolio{
		{"testFirst"},
		{"testSecond"},
	}

	stocks := []models.Stock{
		{"AFLT", 13, 763.03},
		{"GAZP", 130, 21363.00},
	}

	data := make(map[string]interface{})
	data["stocks"] = stocks
	data["portfolios"] = portfolios

	return data
}

func HandleAddPortfolio(w http.ResponseWriter, r *http.Request) {
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

	data := getData()

	err = templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
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

func HandleManager(w http.ResponseWriter, r *http.Request) {
	stocks := []models.Stock{
		{"AFLT", 13, 763.03},
		{"GAZP", 130, 21363.00},
	}

	pieChart, err := services.GetStockPieChart(stocks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := getData()
	data["PieChart"] = base64.StdEncoding.EncodeToString(pieChart)

	cookie, err := r.Cookie("current_portfolio")
	if (err != nil) && (!errors.Is(err, http.ErrNoCookie)) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if cookie == nil || cookie.Value == "" || errors.Is(err, http.ErrNoCookie) {
		handleNoPortfolioChosen(w, r, data)
		return
	}

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
