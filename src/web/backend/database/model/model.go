package model

type Portfolio struct {
	ID   int // PK
	Name string
}

type Asset struct {
	ID    int // PK
	Name  string
	Price float64
}

type PortfolioAsset struct {
	ID          int // PK
	PortfolioID int // FK
	AssetID     int // FK
}

type PortfolioAssetRelationship struct {
	ID               int // PK
	PortfolioAssetID int // FK
	Quantity         int
}

type Index struct {
	ID   int // PK
	Name string
}

type IndexAsset struct {
	ID      int // PK
	IndexID int // FK
	AssetID int // FK
}

type IndexAssetRelationship struct {
	ID           int // PK
	IndexAssetID int // FK
	Fraction     float64
}
