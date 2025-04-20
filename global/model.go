package global

type EditCommentReq struct {
	Id      int64
	Comment string
	Tags    string
}

type GetCombineOrderListReq struct {
	OpenSide  string `json:"openSide" form:"openSide"` // BUY SELL ALL
	Symbol    string `json:"symbol" form:"symbol"`
	DateMin   int64  `json:"dateMin" form:"dateMin"`
	DateMax   int64  `json:"dateMax" form:"dateMax"`
	AmountMin int64  `json:"amountMin" form:"amountMin"`
	AmountMax int64  `json:"amountMax" form:"amountMax"`
	Page      int64  `json:"page" form:"page"`
	PageSize  int64  `json:"pageSize" form:"pageSize"`
	Tags      string `json:"tags" form:"tags"`

	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

type GetPlanListReq struct {
}

type GetPlanListResPlan struct {
	Id           int64  `json:"id"`
	Symbol       string `json:"symbol"`
	OpenPrice    string `json:"openPrice"`
	LossPrice    string `json:"lossPrice"`
	WinPrice     string `json:"winPrice"`
	Notice       int    `json:"notice" `
	AutoTrade    int    `json:"autoTrade"`
	CreateTime   int64  `json:"createTime"`
	PositionSide string `json:"positionSide"`
	// UpdateTime time.Time `json:"updateTime"`
}

type SavePlanReq struct {
	Id           int64  `json:"id" form:"id"`
	Symbol       string `json:"symbol" form:"symbol"`
	OpenPrice    string `json:"openPrice" form:"openPrice"`
	Notice       int    `json:"notice" form:"notice"`
	PositionSide string `json:"positionSide" form:"positionSide"`
	LossPrice    string `json:"lossPrice" form:"lossPrice"`
	WinPrice     string `json:"winPrice" form:"winPrice"`
	AutoTrade    int    `json:"autoTrade" form:"autoTrade"`
}

type DelPlanReq struct {
	Id int64 `json:"id" form:"id"`
}
