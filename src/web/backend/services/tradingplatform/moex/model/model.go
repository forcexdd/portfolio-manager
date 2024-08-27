package model

type AssetData struct {
	TradeDate      string
	BoardID        string
	SecID          string
	TradeTime      string
	CurPrice       float64
	LastPrice      float64
	LegalClose     float64
	TradingSession int
}

type CurrentPricesData struct {
	Data [][]interface{}
}

type IndexData struct {
	IndexID   string
	ShortName string
	From      string
	Till      string
}

type AnalyticsData struct {
	Data [][]interface{}
}

type IndexAssetsData struct {
	IndexID        string
	TradeDate      string
	Ticker         string
	ShortNames     string
	SecIDs         string
	Weight         float64
	TradingSession int
}

type IndexAnalyticsData struct {
	Data [][]interface{}
}
