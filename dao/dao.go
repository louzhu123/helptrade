package dao

import (
	"fmt"
	"main/global"
	"strconv"
	"time"
)

func SelectAccountTradeGroupByOrderId() {
	list := make([]AccountTrade, 0)
	err := global.DB.Model(AccountTrade{}).Group("orderId,symbol,time").
		Select("sum(realizedPnl) as realizedPnl", "symbol", "time").
		Order("time").
		Find(&list).Error
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(list)
	for _, v := range list {
		t := time.UnixMilli(v.Time).Format("2006-01-02 15:04:05")
		fmt.Println(t, v.Symbol, v.RealizedPnl)
	}
}

func GetAllOrder() ([]Order, error) {
	list := make([]Order, 0)
	err := global.DB.Model(Order{}).Find(&list).Order("time asc").Error
	return list, err
}

type TmpCombineOrder struct {
	CurrentPostion  float64
	CurrentCumQuote float64
	Order           CombineOrder
}

// 这个应该卸写在service给定时任务执行
func CombineAccountOrder() []CombineOrder {

	// 仓位状态
	tmpCombineOrder := make(map[string]TmpCombineOrder, 0)

	// 合并数据
	combineOrderList := make([]CombineOrder, 0)

	list, _ := GetAllOrder()

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

		fmt.Println(executedQtyFloat, tmpOrder.CurrentPostion)

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
			combineOrderList = append(combineOrderList, tmpOrder.Order)
			tmpCombineOrder[v.Symbol] = TmpCombineOrder{}
		}
	}

	for _, v := range combineOrderList {
		t := time.UnixMilli(v.StartTime).Format("2006-01-02 15:04:05")
		fmt.Printf("%v %v %v pnl:%.2f \n", t, v.Side, v.Symbol, v.PnL, v.MaxCumQuote)
	}

	return combineOrderList
}

func SaveCombineOrder(list []CombineOrder) {
	err := global.DB.Model(CombineOrder{}).Save(list).Error
	fmt.Println(err)
}

func QueryCombineOrder() ([]CombineOrder, error) {
	var list []CombineOrder
	err := global.DB.Model(&CombineOrder{}).Order("startTime desc").Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func UpdateCombineOrderComment(id int64, comment string) {
	global.DB.Table("combine_order").Where("id", id).Update("comment", comment)
}
