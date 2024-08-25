package moex_models

type AssetData struct {
	TradeDate      string
	BoardId        string
	SecId          string
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
	IndexId   string
	ShortName string
	From      string
	Till      string
}

type AnalyticsData struct {
	Data [][]interface{}
}

type IndexAssetsData struct {
	IndexId        string
	TradeDate      string
	Ticker         string
	ShortNames     string
	SecIds         string
	Weight         float64
	TradingSession int
}

type IndexAnalyticsData struct {
	Data [][]interface{}
}
