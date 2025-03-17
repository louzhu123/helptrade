package global

type EditCommentReq struct {
	Id      int64
	Comment string
}

type GetCombineOrderListReq struct {
	OpenSide  string `json:"openSide" form:"openSide"` // BUY SELL ALL
	Symbol    string `json:"symbol" form:"symbol"`
	DateMin   int64  `json:"dateMin" form:"dateMin"`
	DateMax   int64  `json:"dateMax" form:"dateMax"`
	AmountMin int64  `json:"amountMin" form:"amountMin"`
	AmountMax int64  `json:"amountMax" form:"amountMax"`
}

type GetPlanListReq struct {
}

type GetPlanListResPlan struct {
	Id         int64  `json:"id"`
	Symbol     string `json:"symbol"`
	OpenPrice  string `json:"openPrice"`
	LossPrice  string `json:"lossPrice"`
	WinPrice   string `json:"winPrice"`
	Notice     int    `json:"notice" `
	AutoTrade  int    `json:"autoTrade"`
	CreateTime int64  `json:"createTime"`
	// UpdateTime time.Time `json:"updateTime"`
}

type SavePlanReq struct {
	Id        int64  `json:"id" form:"id"`
	Symbol    string `json:"symbol" form:"symbol"`
	OpenPrice string `json:"openPrice" form:"openPrice"`
}

type DelPlanReq struct {
	Id int64 `json:"id" form:"id"`
}
