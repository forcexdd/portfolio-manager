package chart_drawer

import (
	"bytes"
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/wcharczuk/go-chart/v2"
)

func GetStockPieChart(portfolio *models.Portfolio) ([]byte, error) {
	var values []chart.Value

	for stock, quantity := range portfolio.StocksQuantityMap {
		values = append(values, chart.Value{
			Label: stock.Name,
			Value: stock.Price * float64(quantity),
		})
	}

	pieChart := chart.PieChart{
		Values: values,
	}
	
	var buf bytes.Buffer

	if err := pieChart.Render(chart.SVG, &buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
