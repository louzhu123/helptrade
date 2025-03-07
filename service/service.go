package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"helptrade/dao"
	"helptrade/global"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

func FetchAndCombineOrder() {
	FetchAndSaveAllAccountTrade()

	FetchAndSaveAllOrder()

	list := CombineAccountOrder()
	dao.UpsertCombineOrder(list)
}

func FetchAllAccountTrade() ([]*futures.AccountTrade, error) {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	t := time.Now()
	startTimeUnix := t.Unix()*1000 - 150*24*60*60*1000
	endTimeUnix := startTimeUnix + 7*24*60*60*1000
	maxTime := 1000 // 最多循环1000次应该获取完近90天的数据了
	times := 0
	endFlag := false

	var allData []*futures.AccountTrade
	for {
		t1 := time.Unix(0, startTimeUnix*int64(time.Millisecond))
		t1String := t1.Format("2006-01-02 15:04:05.000")
		t2 := time.Unix(0, endTimeUnix*int64(time.Millisecond))
		t2String := t2.Format("2006-01-02 15:04:05.000")
		fmt.Println("fetch", times, t1String, t2String)

		if endTimeUnix > time.Now().Unix()*1000 {
			endTimeUnix = time.Now().Unix() * 1000
			endFlag = true
		}

		times += 1
		if times > maxTime {
			return nil, errors.New("times too many")
		}
		res, err := futuresClient.NewListAccountTradeService().
			StartTime(startTimeUnix).EndTime(endTimeUnix).Limit(1000).Do(context.Background())
		fmt.Println("len(res)", len(res))
		if err != nil {
			fmt.Println(err)
			if err.Error() == "<APIError> code=-4166, msg=Search window is restricted to recent 90 days only." {
				break
			} else {
				return nil, err
			}
		}

		if len(res) > 0 {
			startTimeUnix = res[len(res)-1].Time + 1
			endTimeUnix = startTimeUnix + 7*24*60*60*1000
		}
		if len(res) == 0 {
			startTimeUnix = endTimeUnix + 1
			endTimeUnix = startTimeUnix + 7*24*60*60*1000
		}
		allData = append(allData, res...)

		if endFlag {
			break
		}
	}

	return allData, nil
}

func FetchAndSaveAllAccountTrade() {
	list, err := FetchAllAccountTrade()
	if err != nil {
		return
	}

	for _, v := range list {
		dao.UpsertAccountTrade(v)
	}
}

func FetchAllOrder() ([]*futures.Order, error) {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	t := time.Now()
	startTimeUnix := t.Unix()*1000 - 90*24*60*60*1000 + 1000*10
	endTimeUnix := startTimeUnix + 7*24*60*60*1000
	maxTime := 1000 // 最多循环1000次应该获取完近90天的数据了
	times := 0
	endFlag := false

	var allData []*futures.Order
	for {
		t1 := time.Unix(0, startTimeUnix*int64(time.Millisecond))
		t1String := t1.Format("2006-01-02 15:04:05.000")
		t2 := time.Unix(0, endTimeUnix*int64(time.Millisecond))
		t2String := t2.Format("2006-01-02 15:04:05.000")
		fmt.Println("fetch", times, t1String, t2String)

		if endTimeUnix > time.Now().Unix()*1000 {
			endTimeUnix = time.Now().Unix() * 1000
			endFlag = true
		}

		times += 1
		if times > maxTime {
			return nil, errors.New("times too many")
		}
		res, err := futuresClient.NewListOrdersService().
			StartTime(startTimeUnix).EndTime(endTimeUnix).Limit(1000).Do(context.Background())
		fmt.Println("len(res)", len(res))
		if err != nil {
			fmt.Println(err)
			if err.Error() == "<APIError> code=-4166, msg=Search window is restricted to recent 90 days only." {
				break
			} else {
				return nil, err
			}
		}

		if len(res) > 0 {
			startTimeUnix = res[len(res)-1].Time + 1
			endTimeUnix = startTimeUnix + 7*24*60*60*1000
		}
		if len(res) == 0 {
			startTimeUnix = endTimeUnix + 1
			endTimeUnix = startTimeUnix + 7*24*60*60*1000
		}
		allData = append(allData, res...)

		if endFlag {
			break
		}
	}

	return allData, nil
}

func FetchAndSaveAllOrder() {
	list, err := FetchAllOrder()
	if err != nil {
		return
	}

	for _, v := range list {
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
