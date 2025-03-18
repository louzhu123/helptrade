package service

import (
	"common/binance"
	"common/binanceWs"
	"common/miaoNotice"
	"fmt"
	"helptrade/dao"
	"strconv"
)

const (
	PlanStatusTodo = 0
	PlanStatusDone = 1
)

func DoPlan() error {
	list, _ := dao.GetAllPlan()

	symbolPrice := binanceWs.HttpGetFutureSymbolCurrentPriceMap()
	for _, plan := range list {
		if plan.Status == PlanStatusDone {
			continue
		}
		if _, ok := symbolPrice[plan.Symbol]; !ok {
			continue
		}

		openPriceFloat, _ := strconv.ParseFloat(plan.OpenPrice, 64)
		currentPrice := symbolPrice[plan.Symbol]

		if plan.PositionSide == "LONG" && currentPrice <= openPriceFloat ||
			plan.PositionSide == "SHORT" && currentPrice >= openPriceFloat {
			// 获取用户信息
			user, _ := dao.GetUserByUserId(plan.UserId)

			// 通知并更新
			req := miaoNotice.SendReq{
				Id:   user.MiaoNoticeId,
				Text: fmt.Sprintf("%v 到达目标价 %v", plan.Symbol, plan.OpenPrice),
			}
			miaoNotice.SendWechat(req)

			client := binance.NewBinanceUsdtContractClient(user.BnApiKey, user.BnApiSecret)
			zhisunPrice, _ := strconv.ParseFloat(plan.LossPrice, 64)
			zhiyingPrice, _ := strconv.ParseFloat(plan.WinPrice, 64)
			if plan.PositionSide == "LONG" {
				client.DuoV3(plan.Symbol, 100, zhisunPrice, zhiyingPrice, currentPrice)
			} else {
				client.KongV3(plan.Symbol, 100, zhisunPrice, zhiyingPrice, currentPrice)
			}

			dao.DonePlan(plan.Id)
		}

	}
	return nil
}
