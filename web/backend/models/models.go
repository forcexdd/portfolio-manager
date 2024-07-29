package models

import "time"

type Stock struct {
	Id    int // PK
	Name  string
	Price float64
	Time  time.Time
}

type UserPortfolio struct {
	Id   int // PK
	Name string
}

type ManyUserPortfoliosWithManyStocks struct {
	Id          int
	PortfolioId int // FK
	StockId     int // FK
	Quantity    int
}

type IndexStock struct {
	Id          int
	NameOfStock string
	fraction    float64
	Time        time.Time
}
