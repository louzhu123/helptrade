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
	nowTime := time.Now()
	for i := 0; i < 30; i++ {
		endTime := nowTime.AddDate(0, 0, -3*i)
		startTime := endTime.AddDate(0, 0, -3)
		res, err := futuresClient.NewListAccountTradeService().
			StartTime(startTime.Unix() * 1000).EndTime(endTime.Unix() * 1000).Limit(1000).Do(context.Background())
		if err != nil {
			if err.Error() == "<APIError> code=-4166, msg=Search window is restricted to recent 90 days only." {
				fmt.Println(err)
			} else {
				// panic(err)
			}
		} else {
			for _, v := range res {
				dao.UpsertAccountTrade(v)
			}
		}
	}
}

func UpdateOrder() {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	nowTime := time.Now()
	for i := 0; i < 30; i++ {
		endTime := nowTime.AddDate(0, 0, -3*i)
		startTime := endTime.AddDate(0, 0, -3)
		res, err := futuresClient.NewListOrdersService().
			StartTime(startTime.Unix() * 1000).EndTime(endTime.Unix() * 1000).Limit(1000).Do(context.Background())
		if err != nil {
			if err.Error() == "<APIError> code=-4166, msg=Search window is restricted to recent 90 days only." {
				fmt.Println(err)
			} else {
				// panic(err)
			}
		} else {
			for _, v := range res {
				dao.UpsertOrder(v)
			}
		}
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
		executedQtyFloat, _ := strconv.ParseFloat(v.ExecutedQty, 64) // 有的标的数量是带一个小数点的，避免浮点数计算问题
		executedQtyFloat100 := executedQtyFloat * 100
		cumQuoteFloat, _ := strconv.ParseFloat(v.CumQuote, 64)
		avgPriceFloat, _ := strconv.ParseFloat(v.AvgPrice, 64)
		if executedQtyFloat100 == 0 { // 无效订单，开了没执行的后面关了的
			continue
		}

		if _, ok := tmpCombineOrder[v.Symbol]; !ok {
			tmpCombineOrder[v.Symbol] = TmpCombineOrder{}
		}
		tmpOrder := tmpCombineOrder[v.Symbol]

		// 手续费
		commission := dao.GetTotalCommissionByOrderId(v.OrderId)
		tmpOrder.Order.Commission += commission

		totalPnl := dao.GetTotalPnlByOrderId(v.OrderId)
		tmpOrder.Order.PnL += totalPnl

		if tmpOrder.CurrentPostion == 0 { // 新开仓
			tmpOrder.Order.StartTime = v.Time
			tmpOrder.Order.PositionSide = v.PositionSide
			tmpOrder.Order.Side = v.Side
			tmpOrder.Order.Symbol = v.Symbol
			tmpOrder.Order.OpenPrice = avgPriceFloat
			tmpOrder.Order.FirstOpenCumQuote = cumQuoteFloat
		} else if v.Side != tmpOrder.Order.Side && executedQtyFloat100-tmpOrder.CurrentPostion == 0 { //结束
			endFlag = true
			tmpOrder.Order.EndTime = v.Time
			tmpOrder.Order.ClosePrice = avgPriceFloat
		}

		if v.Side == tmpOrder.Order.Side {
			tmpOrder.CurrentPostion += executedQtyFloat100
			tmpOrder.CurrentCumQuote += cumQuoteFloat
			tmpOrder.Order.TotalOpenCumQuote += cumQuoteFloat
		} else {
			tmpOrder.CurrentPostion -= executedQtyFloat100
			tmpOrder.CurrentCumQuote -= cumQuoteFloat
			tmpOrder.Order.TotalCloseCumQuote += cumQuoteFloat
		}
		if tmpOrder.CurrentCumQuote > tmpOrder.Order.MaxCumQuote {
			tmpOrder.Order.MaxCumQuote = tmpOrder.CurrentCumQuote
		}

		tmpOrder.OriginOrders = append(tmpOrder.OriginOrders, v)

		tmpCombineOrder[v.Symbol] = tmpOrder

		fmt.Println(tmpOrder.CurrentPostion)

		if endFlag {
			// diff := tmpOrder.Order.TotalCloseCumQuote - tmpOrder.Order.TotalOpenCumQuote
			// tmpOrder.Order.PnL = diff
			// if tmpOrder.Order.Side == "BUY" {
			// 	tmpOrder.Order.PnL = diff
			// } else {
			// 	tmpOrder.Order.PnL = -diff
			// }
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
