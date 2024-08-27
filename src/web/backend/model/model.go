package model

type Portfolio struct {
	Name              string
	AssetsQuantityMap map[*Asset]int
}

type Asset struct {
	Name  string
	Price float64
}

type Index struct {
	Name              string
	AssetsFractionMap map[*Asset]float64
}
