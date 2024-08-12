package models

import "time"

type Portfolio struct {
	Id   int // PK
	Name string
}

type Stock struct {
	Id    int // PK
	Name  string
	Price float64
	Time  time.Time
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

type IndexStock struct {
	Id          int // PK
	NameOfStock string
	fraction    float64
	Time        time.Time
}
