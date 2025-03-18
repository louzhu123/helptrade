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
	"github.com/shopspring/decimal"
)

type Ctx struct {
	User dao.User
}

// func FetchAndCombineOrder() {
// 	users, _ := dao.GetAllUser()
// 	for _, user := range users {
// 		ctx := Ctx{User: user}
// 		FetchAndSaveAllAccountTrade(ctx)

// 		FetchAndSaveAllOrder(ctx)

// 		list := CombineAccountOrder(ctx)
// 		dao.UpsertCombineOrder(list)
// 	}
// }

func FetchAndCombineAccountTrade() {
	users, _ := dao.GetAllUser()
	for _, user := range users {
		fmt.Printf("%+v", user)
		ctx := Ctx{User: user}
		FetchAndSaveAllAccountTrade(ctx)

		FetchAndSaveAllOrder(ctx)

		list := CombineAccountTrade(ctx)
		dao.UpsertCombineOrder(list)
	}
}

func FetchAllAccountTrade(ctx Ctx, startUninxMilli int64) ([]*futures.AccountTrade, error) {
	apiKey := ctx.User.BnApiKey
	apiSecret := ctx.User.BnApiSecret
	futuresClient := binance.NewFuturesClient(apiKey, apiSecret)
	t := time.Now()
	startTimeUnix := t.Unix()*1000 - 500*24*60*60*1000 // 默认500天前开始 TODO 其他用户应该还有甚至几年的
	if startUninxMilli != 0 {
		startTimeUnix = startUninxMilli - 1000*200
	}
	endTimeUnix := startTimeUnix + 7*24*60*60*1000
	maxTime := 1000 // 最多循环1000次应该获取完近90天的数据了
	times := 0
	endFlag := false
	lastTimeStamp := int64(0)

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
			if lastTimeStamp == res[len(res)-1].Time {
				// 死循环了，跳过
				startTimeUnix = endTimeUnix - 1000*100
				endTimeUnix = startTimeUnix + 7*24*60*60*1000
			} else {
				startTimeUnix = res[len(res)-1].Time - 1000*100
				endTimeUnix = startTimeUnix + 7*24*60*60*1000
			}

			lastTimeStamp = res[len(res)-1].Time
		}
		if len(res) == 0 { // 可能这些数据都在1秒钟发生就会导致会死循环
			startTimeUnix = endTimeUnix - 1000*100
			endTimeUnix = startTimeUnix + 7*24*60*60*1000
		}
		allData = append(allData, res...)

		if endFlag {
			break
		}
	}

	return allData, nil
}

func FetchAndSaveAllAccountTrade(ctx Ctx) {
	// 获取最新一条数据
	latest, _ := dao.GetLastestAccountTradeByUserId(ctx.User.Id)
	startTime := latest.Time

	list, err := FetchAllAccountTrade(ctx, startTime) // 数据会重复，为了能获取完整，时间窗口会有重叠
	if err != nil {
		return
	}

	for _, v := range list {
		data := dao.AccountTrade{}
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &data)
		data.UserId = ctx.User.Id
		dao.UpsertAccountTrade(data)
	}
}

func FetchAllOrder(ctx Ctx, startUninxMilli int64) ([]*futures.Order, error) {
	apiKey := ctx.User.BnApiKey
	apiSecret := ctx.User.BnApiSecret
	futuresClient := binance.NewFuturesClient(apiKey, apiSecret)
	t := time.Now()
	startTimeUnix := t.Unix()*1000 - 90*24*60*60*1000 + 1000*10 // 加个10s，访问过程中已经过了一秒的90天窗口
	if startUninxMilli != 0 {
		startTimeUnix = startTimeUnix - 1000*200
	}

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

func FetchAndSaveAllOrder(ctx Ctx) {
	lastest, _ := dao.GetLastestOrderByUserId(ctx.User.Id)
	startTime := lastest.Time
	list, err := FetchAllOrder(ctx, startTime)
	if err != nil {
		return
	}

	for _, v := range list {
		data := dao.Order{}
		b, _ := json.Marshal(v)
		json.Unmarshal(b, &data)
		data.UserId = ctx.User.Id
		dao.UpsertOrder(data)
	}
}

type TmpCombineOrder struct {
	CurrentPostion  float64
	CurrentCumQuote float64
	OriginOrders    []dao.Order
	Order           dao.CombineOrder
}

func CombineAccountOrder(ctx Ctx) []dao.CombineOrder {

	// 仓位状态
	tmpCombineOrder := make(map[string]TmpCombineOrder, 0)

	// 合并数据
	combineOrderList := make([]dao.CombineOrder, 0)

	list, _ := dao.GetAllOrderByUserId(ctx.User.Id)

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
			tmpOrder.Order.UserId = ctx.User.Id
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

type TmpCombineAccountTrade struct {
	CurrentPostion  decimal.Decimal
	CurrentCumQuote float64
	OriginOrders    []dao.AccountTrade
	Order           dao.CombineOrder
}

func CombineAccountTrade(ctx Ctx) []dao.CombineOrder {

	// 仓位状态
	tmpCombineOrder := make(map[string]TmpCombineAccountTrade, 0)

	// 合并数据
	combineOrderList := make([]dao.CombineOrder, 0)

	list, _ := dao.GetAccountTradeByUserId(ctx.User.Id)
	// fmt.Println(len(list))

	for _, v := range list {
		endFlag := false
		fmt.Println("当前accountrade", v.Qty)
		// executedQtyFloat, _ := strconv.ParseFloat(v.Qty, 64) // 有的标的数量是带一个小数点的，避免浮点数计算问
		// executedQtyFloat100 := executedQtyFloat * 100
		executedQtyDecimal, _ := decimal.NewFromString(v.Qty)
		executedQtyFloat, _ := executedQtyDecimal.Float64()

		commissionFloat, _ := strconv.ParseFloat(v.Commission, 64)
		pnl, _ := strconv.ParseFloat(v.RealizedPnl, 64)

		cumQuoteFloat, _ := strconv.ParseFloat(v.QuoteQty, 64)
		avgPriceFloat, _ := strconv.ParseFloat(v.Price, 64)
		if executedQtyFloat == 0 { // 无效订单，开了没执行的后面关了的
			continue
		}

		if _, ok := tmpCombineOrder[v.Symbol]; !ok {
			tmpCombineOrder[v.Symbol] = TmpCombineAccountTrade{}
		}
		tmpOrder := tmpCombineOrder[v.Symbol]

		currentPostionFloat, _ := tmpOrder.CurrentPostion.Float64()

		// fmt.Println(tmpOrder.CurrentPostion)

		tmpOrder.Order.Commission += commissionFloat
		tmpOrder.Order.PnL += pnl

		subResult, _ := executedQtyDecimal.Sub(tmpOrder.CurrentPostion).Float64()

		if currentPostionFloat == 0 { // 新开仓
			fmt.Println("新开仓")
			tmpOrder.Order.StartTime = v.Time
			tmpOrder.Order.PositionSide = v.PositionSide
			tmpOrder.Order.Side = v.Side
			tmpOrder.Order.Symbol = v.Symbol
			tmpOrder.Order.OpenPrice = avgPriceFloat
			tmpOrder.Order.FirstOpenCumQuote = cumQuoteFloat
		} else if v.Side != tmpOrder.Order.Side && subResult == 0 { //结束
			endFlag = true
			tmpOrder.Order.EndTime = v.Time
			tmpOrder.Order.ClosePrice = avgPriceFloat
		}

		if v.Side == tmpOrder.Order.Side {
			tmpOrder.CurrentPostion = tmpOrder.CurrentPostion.Add(executedQtyDecimal)
			tmpOrder.CurrentCumQuote += cumQuoteFloat
			tmpOrder.Order.TotalOpenCumQuote += cumQuoteFloat
		} else {
			tmpOrder.CurrentPostion = tmpOrder.CurrentPostion.Sub(executedQtyDecimal)
			tmpOrder.CurrentCumQuote -= cumQuoteFloat
			tmpOrder.Order.TotalCloseCumQuote += cumQuoteFloat
		}

		fmt.Println("当前仓位", tmpOrder.CurrentPostion)
		if tmpOrder.CurrentCumQuote > tmpOrder.Order.MaxCumQuote {
			tmpOrder.Order.MaxCumQuote = tmpOrder.CurrentCumQuote
		}

		tmpOrder.OriginOrders = append(tmpOrder.OriginOrders, v)

		tmpCombineOrder[v.Symbol] = tmpOrder

		if endFlag {
			fmt.Println("平仓所有\n\n\n\n")
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
			tmpOrder.Order.UserId = ctx.User.Id
			combineOrderList = append(combineOrderList, tmpOrder.Order)
			tmpCombineOrder[v.Symbol] = TmpCombineAccountTrade{}
		}
	}

	// for _, v := range combineOrderList {
	// 	t := time.UnixMilli(v.StartTime).Format("2006-01-02 15:04:05")
	// 	fmt.Printf("%v %v %v pnl:%.2f \n", t, v.Side, v.Symbol, v.PnL)
	// }

	return combineOrderList
}

func GetPlanList() ([]global.GetPlanListResPlan, error) {
	var res []global.GetPlanListResPlan
	list, err := dao.GetAllPlan()
	for _, item := range list {
		tmp := global.GetPlanListResPlan{
			Id:           item.Id,
			Symbol:       item.Symbol,
			OpenPrice:    item.OpenPrice,
			LossPrice:    item.LossPrice,
			WinPrice:     item.LossPrice,
			Notice:       item.Notice,
			AutoTrade:    item.AutoTrade,
			PositionSide: item.PositionSide,
			CreateTime:   item.CreateTime.Unix(),
		}
		res = append(res, tmp)
	}
	return res, err
}

func SavePlan(req global.SavePlanReq) error {
	if req.Id != 0 { // 更新
		data, err := dao.GetPlanById(req.Id)
		if err != nil {
			return err
		}
		if data.Id == 0 {
			return errors.New("id 不存在")
		}
		data.OpenPrice = req.OpenPrice
		data.Symbol = req.Symbol
		data.Notice = req.Notice
		data.PositionSide = req.PositionSide
		dao.SavePlan(data)
	} else { // 新增
		data := dao.Plan{
			Symbol:       req.Symbol,
			OpenPrice:    req.OpenPrice,
			Notice:       req.Notice,
			PositionSide: req.PositionSide,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := dao.CreatePlan(&data)
		return err
	}

	return nil
}

func DelPlan(req global.DelPlanReq) error {
	dao.DelPlan(req.Id)
	return nil
}

func GetUserByToken(token string) (dao.User, error) {
	return dao.GetUserByToken(token)
}
