package dao

import "time"

type Order struct {
	Symbol                  string `json:"symbol" gorm:"column:symbol"`
	OrderId                 int64  `json:"orderId" gorm:"column:orderId"`
	ClientOrderId           string `json:"clientOrderId" gorm:"column:clientOrderId"`
	Price                   string `json:"price" gorm:"column:price"`
	ReduceOnly              bool   `json:"reduceOnly" gorm:"column:reduceOnly"`
	OrigQty                 string `json:"origQty" gorm:"column:origQty"`
	ExecutedQty             string `json:"executedQty" gorm:"column:executedQty"`
	CumQty                  string `json:"cumQty" gorm:"column:cumQty"`
	CumQuote                string `json:"cumQuote" gorm:"column:cumQuote"`
	Status                  string `json:"status" gorm:"column:status"`
	TimeInForce             string `json:"timeInForce" gorm:"column:timeInForce"`
	Type                    string `json:"type" gorm:"column:type"`
	Side                    string `json:"side" gorm:"column:side"`
	StopPrice               string `json:"stopPrice" gorm:"column:stopPrice"`
	Time                    int64  `json:"time" gorm:"column:time"`
	UpdateTime              int64  `json:"updateTime" gorm:"column:updateTime"`
	WorkingType             string `json:"workingType" gorm:"column:workingType"`
	ActivatePrice           string `json:"activatePrice" gorm:"column:activatePrice"`
	PriceRate               string `json:"priceRate" gorm:"column:priceRate"`
	AvgPrice                string `json:"avgPrice" gorm:"column:avgPrice"`
	OrigType                string `json:"origType" gorm:"column:origType"`
	PositionSide            string `json:"positionSide" gorm:"column:positionSide"`
	PriceProtect            bool   `json:"priceProtect" gorm:"column:priceProtect"`
	ClosePosition           bool   `json:"closePosition" gorm:"column:closePosition"`
	PriceMatch              string `json:"priceMatch" gorm:"column:priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode" gorm:"column:selfTradePreventionMode"`
	GoodTillDate            int    `json:"goodTillDate" gorm:"column:goodTillDate"`
	UserId                  int    `json:"userId" gorm:"column:userId"`
}

func (Order) TableName() string {
	return "order"
}

type AccountTrade struct {
	Buyer           bool   `json:"buyer" gorm:"column:buyer;type:tinyint(1);not_null"`
	Commission      string `json:"commission" gorm:"column:commission;type:varchar(255);not_null"`
	CommissionAsset string `json:"commissionAsset" gorm:"column:commissionAsset;type:varchar(255);not_null"`
	ID              int    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	Maker           bool   `json:"maker" gorm:"column:maker;type:tinyint(1);not_null"`
	OrderID         int    `json:"orderId" gorm:"column:orderId;type:bigint;not_null"`
	Price           string `json:"price" gorm:"column:price;type:varchar(255);not_null"`
	Qty             string `json:"qty" gorm:"column:qty;type:varchar(255);not_null"`
	QuoteQty        string `json:"quoteQty" gorm:"column:quoteQty;type:varchar(255);not_null"`
	RealizedPnl     string `json:"realizedPnl" gorm:"column:realizedPnl;type:varchar(255);not_null"`
	Side            string `json:"side" gorm:"column:side;type:varchar(255);not_null"`
	PositionSide    string `json:"positionSide" gorm:"column:positionSide;type:varchar(255);not_null"`
	Symbol          string `json:"symbol" gorm:"column:symbol;type:varchar(255);not_null"`
	Time            int64  `json:"time" gorm:"column:time;type:bigint;not_null"`
	UserId          int    `json:"userId" gorm:"column:userId"`
}

func (AccountTrade) TableName() string {
	return "account_trade"
}

type CombineOrder struct {
	Id                 int64   `json:"id" gorm:"column:id"`
	Symbol             string  `json:"symbol"`
	PnL                float64 `json:"pnl" gorm:"column:pnl"`
	StartTime          int64   `json:"startTime" gorm:"column:startTime"`
	EndTime            int64   `json:"endTime" gorm:"column:endTime"`
	PositionSide       string  `json:"positionSide" gorm:"column:positionSide"`
	Side               string  `json:"side"`
	FirstOpenCumQuote  float64 `json:"firstOpenCumQuote" gorm:"column:firstOpenCumQuote"`
	TotalOpenCumQuote  float64 `json:"totalOpenCumQuote" gorm:"column:totalOpenCumQuote"`
	TotalCloseCumQuote float64 `json:"totalCloseCumQuote" gorm:"column:totalCloseCumQuote"`
	MaxCumQuote        float64 `json:"maxCumQuote" gorm:"column:maxCumQuote"`
	Comment            string  `json:"comment" gorm:"column:comment"`
	OpenPrice          float64 `json:"openPrice" gorm:"column:openPrice"`
	ClosePrice         float64 `json:"closePrice" gorm:"column:closePrice"`
	OriginOrders       string  `json:"originOrders" gorm:"column:originOrders"`
	Commission         float64 `json:"commission" gorm:"column:commission"`
	UserId             int     `json:"userId" gorm:"column:userId"`
}

func (CombineOrder) TableName() string {
	return "combine_order"
}

type Plan struct {
	Id           int64     `json:"id" gorm:"column:id"`
	Symbol       string    `json:"symbol"`
	OpenPrice    string    `json:"openPrice" gorm:"column:openPrice"`
	LossPrice    string    `json:"lossPrice" gorm:"column:lossPrice"`
	WinPrice     string    `json:"winPrice" gorm:"column:winPrice"`
	Notice       int       `json:"notice" gorm:"column:notice"`
	AutoTrade    int       `json:"autoTrade" gorm:"column:autoTrade"`
	CreateTime   time.Time `json:"createTime" gorm:"column:createTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:updateTime"`
	Status       int       `json:"status" gorm:"column:status"`
	PositionSide string    `json:"positionSide" gorm:"column:positionSide"`
	UserId       int       `json:"userId" gorm:"column:userId"`
}

func (Plan) TableName() string {
	return "plan"
}

type User struct {
	Id           int    `json:"id" gorm:"column:id"`
	Token        string `json:"token" gorm:"column:token"`
	BnApiKey     string `json:"bnApiKey" gorm:"column:bnApiKey"`
	BnApiSecret  string `json:"bnApiSecret" gorm:"column:bnApiSecret"`
	MiaoNoticeId string `json:"miaoNoticeId" gorm:"column:miaoNoticeId"`
}

func (User) TableName() string {
	return "user"
}
