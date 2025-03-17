package service

import (
	"common/binanceWs"
	"common/miaoNotice"
	"fmt"
	"helptrade/dao"
	"helptrade/global"
	"strconv"
)

const (
	PlanStatusTodo = 0
	PlanStatusDone = 1
)

func DoPlan() error {
	list, _ := dao.GetAllPlan()

	symbolPrice := binanceWs.HttpGetFutureSymbolCurrentPriceMap()
	for _, item := range list {
		if item.Status == PlanStatusDone {
			continue
		}
		if _, ok := symbolPrice[item.Symbol]; !ok {
			continue
		}

		openPriceFloat, _ := strconv.ParseFloat(item.OpenPrice, 64)
		currentPrice := symbolPrice[item.Symbol]

		if item.PositionSide == "LONG" && currentPrice <= openPriceFloat ||
			item.PositionSide == "SHORT" && currentPrice >= openPriceFloat {
			// 通知并更新
			req := miaoNotice.SendReq{
				Id:   global.Cfg.MiaoNoticeId,
				Text: fmt.Sprintf("%v 到达目标价 %v", item.Symbol, item.OpenPrice),
			}
			miaoNotice.SendWechat(req)

			dao.DonePlan(item.Id)
		}

	}
	return nil
}
