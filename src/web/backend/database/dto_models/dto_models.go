package dto_models

type Portfolio struct {
	Id   int // PK
	Name string
}

type Asset struct {
	Id    int // PK
	Name  string
	Price float64
}

type PortfolioAsset struct {
	Id          int // PK
	PortfolioId int // FK
	AssetId     int // FK
}

type PortfolioAssetRelationship struct {
	Id               int // PK
	PortfolioAssetId int // FK
	Quantity         int
}

type Index struct {
	Id   int // PK
	Name string
}

type IndexAsset struct {
	Id      int // PK
	IndexId int // FK
	AssetId int // FK
}

type IndexAssetRelationship struct {
	Id           int // PK
	IndexAssetId int // FK
	Fraction     float64
}
