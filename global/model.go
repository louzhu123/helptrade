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
