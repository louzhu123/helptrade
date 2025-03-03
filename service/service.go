package service

import (
	"context"
	"fmt"
	"helptrade/dao"
	"helptrade/global"

	"github.com/adshao/go-binance/v2"
)

func FetchAndCombineOrder() {
	UpdateAccountTrade()

	UpdateOrder()

	list := dao.CombineAccountOrder()
	dao.UpsertCombineOrder(list)
}

func UpdateAccountTrade() {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	res, err := futuresClient.NewListAccountTradeService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		dao.UpsertAccountTrade(v)
	}
}

func UpdateOrder() {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	res, err := futuresClient.NewListOrdersService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		dao.UpsertOrder(v)
	}
}
