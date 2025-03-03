package service

import (
	"context"
	"encoding/json"
	"fmt"
	"helptrade/dao"
	"helptrade/global"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

func FetchAndCombineOrder() {
	UpdateAccountTrade()

	UpdateOrder()

	list := CombineAccountOrder()
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

type TmpCombineOrder struct {
	CurrentPostion  float64
	CurrentCumQuote float64
	OriginOrders    []dao.Order
	Order           dao.CombineOrder
}

func CombineAccountOrder() []dao.CombineOrder {

	// 仓位状态
	tmpCombineOrder := make(map[string]TmpCombineOrder, 0)

	// 合并数据
	combineOrderList := make([]dao.CombineOrder, 0)

	list, _ := dao.GetAllOrder()

	for _, v := range list {
		endFlag := false
		executedQtyFloat, _ := strconv.ParseFloat(v.ExecutedQty, 64)
		cumQuoteFloat, _ := strconv.ParseFloat(v.CumQuote, 64)
		avgPriceFloat, _ := strconv.ParseFloat(v.AvgPrice, 64)
		if executedQtyFloat == 0 { // 无效订单，开了没执行的后面关了的
			continue
		}

		if _, ok := tmpCombineOrder[v.Symbol]; !ok {
			tmpCombineOrder[v.Symbol] = TmpCombineOrder{}
		}
		tmpOrder := tmpCombineOrder[v.Symbol]

		// 手续费
		commission := dao.GetTotalCommissionByOrderId(v.OrderId)
		tmpOrder.Order.Commission += commission

		if tmpOrder.CurrentPostion == 0 { // 新开仓
			tmpOrder.Order.StartTime = v.Time
			tmpOrder.Order.PositionSide = v.PositionSide
			tmpOrder.Order.Side = v.Side
			tmpOrder.Order.Symbol = v.Symbol
			tmpOrder.Order.OpenPrice = avgPriceFloat
			tmpOrder.Order.FirstOpenCumQuote = cumQuoteFloat
		} else if v.Side != tmpOrder.Order.Side && executedQtyFloat-tmpOrder.CurrentPostion == 0 { //结束
			endFlag = true
			tmpOrder.Order.EndTime = v.Time
			tmpOrder.Order.ClosePrice = avgPriceFloat
		}

		if v.Side == tmpOrder.Order.Side {
			tmpOrder.CurrentPostion += executedQtyFloat
			tmpOrder.CurrentCumQuote += cumQuoteFloat
			tmpOrder.Order.TotalOpenCumQuote += cumQuoteFloat
		} else {
			tmpOrder.CurrentPostion -= executedQtyFloat
			tmpOrder.CurrentCumQuote -= cumQuoteFloat
			tmpOrder.Order.TotalCloseCumQuote += cumQuoteFloat
		}
		if tmpOrder.CurrentCumQuote > tmpOrder.Order.MaxCumQuote {
			tmpOrder.Order.MaxCumQuote = tmpOrder.CurrentCumQuote
		}

		tmpOrder.OriginOrders = append(tmpOrder.OriginOrders, v)

		tmpCombineOrder[v.Symbol] = tmpOrder

		if endFlag {
			diff := tmpOrder.Order.TotalCloseCumQuote - tmpOrder.Order.TotalOpenCumQuote
			tmpOrder.Order.PnL = diff
			if tmpOrder.Order.Side == "BUY" {
				tmpOrder.Order.PnL = diff
			} else {
				tmpOrder.Order.PnL = -diff
			}
			if tmpOrder.Order.MaxCumQuote < tmpOrder.Order.TotalCloseCumQuote {
				tmpOrder.Order.MaxCumQuote = tmpOrder.Order.TotalCloseCumQuote
			}

			originOrdersStr, _ := json.Marshal(tmpOrder.OriginOrders)
			tmpOrder.Order.OriginOrders = string(originOrdersStr)
			combineOrderList = append(combineOrderList, tmpOrder.Order)
			tmpCombineOrder[v.Symbol] = TmpCombineOrder{}
		}
	}

	for _, v := range combineOrderList {
		t := time.UnixMilli(v.StartTime).Format("2006-01-02 15:04:05")
		fmt.Printf("%v %v %v pnl:%.2f \n", t, v.Side, v.Symbol, v.PnL)
	}

	return combineOrderList
}
