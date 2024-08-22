package moex_models

type StockData struct {
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
