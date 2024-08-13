package models

type Portfolio struct {
	Id   int // PK
	Name string
}

type Stock struct {
	Id    int // PK
	Name  string
	Price float64
}

type PortfolioStock struct {
	Id          int // PK
	PortfolioId int // FK
	StockId     int // FK
}

type PortfolioStockRelationship struct {
	Id               int // PK
	PortfolioStockId int // FK
	Quantity         int
}

type Index struct {
	Id   int // PK
	Name string
}

type IndexStock struct {
	Id      int // PK
	IndexId int // FK
	StockId int // FK
}

type IndexStockRelationship struct {
	Id           int // PK
	IndexStockId int // FK
	fraction     float64
}
