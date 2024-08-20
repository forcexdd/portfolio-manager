package models

type Portfolio struct {
	Name              string
	StocksQuantityMap map[*Stock]int
}

type Stock struct {
	Name  string
	Price float64
}

type Index struct {
	Name              string
	StocksFractionMap map[*Stock]float64
}
