package domain

type Portfolio struct {
	ID   int64
	Name string
}

type Asset struct {
	ID    int64
	Name  string
	Price float64
}

type PortfolioAsset struct {
	Asset    Asset
	Quantity int64
	LotSize  int64
}

type Index struct {
	ID   int64
	Name string
}

type IndexAsset struct {
	Asset    Asset
	Fraction float64
}

type AnalysisResult struct {
	TotalValue float64
	Assets     []AssetDiff
}

type AssetDiff struct {
	AssetID         int64
	Name            string
	Price           float64
	CurrentQuantity int64
	CurrentFraction float64
	IndexFraction   float64
	FractionDiff    float64
	TargetQuantity  int64
	Difference      int64
	LotSize         int64
	Action          string
}

type ImportExportAsset struct {
	Name     string
	Quantity int64
	LotSize  int64
}

type PortfolioProfile struct {
	Name   string
	Assets []ImportExportAsset
}
