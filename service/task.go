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
	fmt.Println("doPlan")
	list, _ := dao.GetAllPlan()

	// todo 暂时一个人用
	tmpUser, _ := dao.GetUserByUserId(1)
	tmpClient := binance.NewBinanceUsdtContractClient(tmpUser.BnApiKey, tmpUser.BnApiSecret)
	account, err := tmpClient.GetAccount()
	if err != nil {
		return nil
	}
	for _, item := range account.Positions {
		amount, _ := strconv.ParseFloat(item.PositionAmt, 64)
		if amount > 0 || amount < 0 {
			return nil
		}
	}

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

		fmt.Println("currentPrice", currentPrice, "openPriceFloat", openPriceFloat)

		if plan.PositionSide == "LONG" && currentPrice <= openPriceFloat ||
			plan.PositionSide == "SHORT" && currentPrice >= openPriceFloat {
			// 获取用户信息
			user, _ := dao.GetUserByUserId(plan.UserId)

			// 通知并更新
			if plan.Notice == 1 {
				req := miaoNotice.SendReq{
					Id:   user.MiaoNoticeId,
					Text: fmt.Sprintf("%v 到达目标价 %v", plan.Symbol, plan.OpenPrice),
				}
				miaoNotice.SendWechat(req)
			}

			if plan.AutoTrade == 1 {
				client := binance.NewBinanceUsdtContractClient(user.BnApiKey, user.BnApiSecret)
				zhisunPrice, _ := strconv.ParseFloat(plan.LossPrice, 64)
				zhiyingPrice, _ := strconv.ParseFloat(plan.WinPrice, 64)
				if plan.PositionSide == "LONG" {
					client.DuoV3(plan.Symbol, 300, zhisunPrice, zhiyingPrice, currentPrice)
				} else {
					client.KongV3(plan.Symbol, 300, zhisunPrice, zhiyingPrice, currentPrice)
				}
			}

			dao.DonePlan(plan.Id)
		}

	}
	return nil
}
