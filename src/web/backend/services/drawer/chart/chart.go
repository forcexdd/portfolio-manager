package chart

import (
	"bytes"

	"github.com/forcexdd/portfoliomanager/src/model"
	"github.com/wcharczuk/go-chart/v2"
)

func GetAssetPieChart(portfolio *model.Portfolio) ([]byte, error) {
	var values []chart.Value

	for asset, quantity := range portfolio.AssetsQuantityMap {
		values = append(values, chart.Value{
			Label: asset.Name,
			Value: asset.Price * float64(quantity),
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
