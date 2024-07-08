package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMysql(dsn string) {
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open mysql failed,", err)
	}
	DB = d
}

type TradeRecord struct {
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
}

func (TradeRecord) TableName() string {
	return "trade_record"
}

func save() {

}
