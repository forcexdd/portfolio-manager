package services

import (
	"github.com/forcexdd/StockPortfolioManager/src/web/backend/models"
	"github.com/vicanso/go-charts"
)

func GetStockPieChart(stocks []models.Stock) ([]byte, error) {
	var values []float64
	var names []string

	for _, stock := range stocks {
		values = append(values, float64(stock.Quantity))
		names = append(names, stock.Name)
	}
	pieChart, err := charts.PieRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text: "PieChart of stocks",
			Left: charts.PositionCenter,
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			Orient: charts.OrientVertical,
			Data:   names,
			Left:   charts.PositionLeft,
		}),
	)
	if err != nil {
		return nil, err
	}
	buf, err := pieChart.Bytes()
	if err != nil {
		return nil, err
	}

	return buf, nil
}
