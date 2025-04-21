package dao

import (
	"fmt"
	"helptrade/global"

	"gorm.io/gorm"
)

func UpdateCombineOrderComment(userId, id int64, comment, tags string) {
	err := global.DB.Table("combine_order").Where("userId", userId).Where("id", id).Updates(
		map[string]interface{}{
			"comment": comment,
			"tags":    tags,
		},
	).Error
	fmt.Println("err", err)
}

func QueryCombineOrder(userId int, req global.GetCombineOrderListReq) ([]CombineOrder, int64, error) {
	where := global.DB.Model(&CombineOrder{})
	if req.OpenSide == "BUY" {
		where.Where("positionSide", "LONG")
	} else if req.OpenSide == "SELL" {
		where.Where("positionSide", "SHORT")
	}

	if req.Symbol != "ALL" && req.Symbol != "" {
		where.Where("symbol", req.Symbol)
	}

	if req.DateMax != 0 {
		where.Where("startTime <= ?", req.DateMax*1000)
	}
	if req.DateMin != 0 {
		where.Where("startTIme >= ?", req.DateMin*1000)
	}
	if req.AmountMax != 0 {
		where.Where("maxCumQuote <= ?", req.AmountMax)
	}
	if req.AmountMin != 0 {
		where.Where("maxCumQuote >= ?", req.AmountMin)
	}
	if req.Tags != "" {
		where.Where("tags LIKE ?", "%"+req.Tags+"%")
	}

	where.Where("userId", userId)

	var count int64
	where.Count(&count)

	if req.Page != 0 && req.PageSize != 0 {
		offset := (req.Page - 1) * req.PageSize
		where.Offset(int(offset)).Limit(int(req.PageSize))
	}

	var list []CombineOrder
	if req.SortOrder != "" {
		where = where.Order(fmt.Sprintf("%s %s", req.SortBy, req.SortOrder))
	} else {
		where = where.Order("startTime desc")
	}
	err := where.Find(&list).Error
	if err != nil {
		return list, count, err
	}
	fmt.Println("list", list)
	return list, count, nil
}

func GetCombineOrderStatis(userId int, req global.GetCombineOrderListReq) (CombineOrderStatis, error) {
	where := global.DB.Model(&CombineOrder{})

	// Select("sum(commission) as totalCommission").
	// Select("sum(pnl) as totalPnl")
	// Select("avg(endTime - startTime) as avgTakeTime").
	// Select("avg(firstOpenCumQuote) as avgFirstOpenCumQuote").
	// Select("avg(maxCumQuote) as avgMaxCumQuote").
	// Select("avg(IF(pnl > 0, pnl, NULL)) as avgWin").
	// Select("avg(IF(pnl < 0, pnl, NULL)) as avgLoss").
	// Select("avg(IF((pnl+commission) > 0, (pnl+commission), NULL)) as avgWinWithCommission").
	// Select("avg(IF((pnl+commission) < 0, (pnl+commission), NULL)) as avgLossWithCommission").
	// Select("count(IF((pnl+commission) > 0, 1, NULL)) as winTimes").
	// Select("count(IF((pnl+commission) < 0, 1, NULL)) as lossTimes").

	if req.OpenSide == "BUY" {
		where.Where("positionSide", "LONG")
	} else if req.OpenSide == "SELL" {
		where.Where("positionSide", "SHORT")
	}

	if req.Symbol != "ALL" && req.Symbol != "" {
		where.Where("symbol", req.Symbol)
	}

	if req.DateMax != 0 {
		where.Where("startTime <= ?", req.DateMax*1000)
	}
	if req.DateMin != 0 {
		where.Where("startTIme >= ?", req.DateMin*1000)
	}
	if req.AmountMax != 0 {
		where.Where("maxCumQuote <= ?", req.AmountMax)
	}
	if req.AmountMin != 0 {
		where.Where("maxCumQuote >= ?", req.AmountMin)
	}
	if req.Tags != "" {
		where.Where("tags LIKE ?", "%"+req.Tags+"%")
	}

	where.Where("userId", userId)

	var data CombineOrderStatis
	selectFileds := []string{
		"sum(commission) as totalCommission",
		"sum(pnl) as totalPnl",
		"sum(pnl-commission) as totalPnlWithCommission",
		"sum(IF(pnl > 0, pnl, NULL)) as totalWinPnl",
		"sum(IF(pnl < 0, pnl, NULL)) as totalLossPnl",
		"count(IF((pnl-commission) > 0, 1, NULL)) as winTimes",
		"count(IF((pnl-commission) < 0, 1, NULL)) as lossTimes",
		"avg(endTime - startTime) as avgTakeTime",
		"avg(IF((pnl-commission) > 0, pnl-commission, NULL)) as avgWinWithCommission",
		"avg(IF((pnl-commission) < 0, pnl-commission, NULL)) as avgLossWithCommission",
		"avg(IF(pnl > 0, pnl, NULL)) as avgWin",
		"avg(IF(pnl < 0, pnl, NULL)) as avgLoss",
		"avg(FirstOpenCumQuote) as avgFirstOpenCumQuote",
		"avg(maxCumQuote) as avgMaxCumQuote",
	}
	where = where.Select(selectFileds)

	err := where.First(&data).Error

	if err != nil {
		return data, err
	}
	return data, nil
}

func GetTotalCommissionByOrderId(orderId int64) float64 {
	var data []float64
	global.DB.Model(AccountTrade{}).Select("sum(commission) as commission").
		Where("orderId", orderId).Pluck("commission", &data)

	if len(data) > 0 {
		return data[0]
	} else {
		return 0
	}
}

func GetTotalPnlByOrderId(orderId int64) float64 {
	var data []float64
	global.DB.Model(AccountTrade{}).Select("sum(realizedPnl) as realizedPnl").
		Where("orderId", orderId).Pluck("realizedPnl", &data)

	if len(data) > 0 {
		return data[0]
	} else {
		return 0
	}
}

func GetAllOrder() ([]Order, error) {
	list := make([]Order, 0)
	err := global.DB.Model(Order{}).Order("time asc").Find(&list).Error
	return list, err
}

func GetAllOrderByUserId(userId int) ([]Order, error) {
	list := make([]Order, 0)
	err := global.DB.Model(Order{}).Where("userId", userId).Order("time asc").Find(&list).Error
	return list, err
}

func GetAllAccountTrade() ([]AccountTrade, error) {
	list := make([]AccountTrade, 0)
	err := global.DB.Model(AccountTrade{}).
		// Order("time asc").Find(&list).Error
		Select("min(userId) as userId,orderId,min(time) as time,sum(commission) as commission,sum(qty) as qty,sum(quoteQty) as quoteQty,sum(realizedPnl) as realizedPnl,avg(price) as price,MIN(symbol) as symbol,MIN(side) as side,MIN(positionSide) as positionSide").
		Group("orderId").Order("time asc").Find(&list).Error
	return list, err
}

func GetLastestAccountTradeByUserId(userId int) (AccountTrade, error) {
	var data AccountTrade
	err := global.DB.Model(AccountTrade{}).Where("userId", userId).Order("time desc").First(&data).Error
	return data, err
}

func GetLastestOrderByUserId(userId int) (Order, error) {
	var data Order
	err := global.DB.Model(Order{}).Where("userId", userId).Order("time desc").First(&data).Error
	return data, err
}

func GetAccountTradeByUserId(userId int) ([]AccountTrade, error) {
	list := make([]AccountTrade, 0)
	err := global.DB.Model(AccountTrade{}).
		// Order("time asc").Find(&list).Error
		Where("userId", userId).
		Select("min(userId) as userId,orderId,min(time) as time,sum(commission) as commission,sum(qty) as qty,sum(quoteQty) as quoteQty,sum(realizedPnl) as realizedPnl,avg(price) as price,MIN(symbol) as symbol,MIN(side) as side,MIN(positionSide) as positionSide").
		Group("orderId").Order("time asc").Find(&list).Error
	return list, err
}

func SaveCombineOrder(list []CombineOrder) error {
	err := global.DB.Model(CombineOrder{}).Save(list).Error
	return err
}

func UpsertCombineOrder(list []CombineOrder) {
	for _, v := range list {
		m := CombineOrder{}
		// 根据开仓时间和标的即可标识
		where := global.DB.Model(CombineOrder{}).Where("startTime = ?", v.StartTime).Where("symbol = ?", v.Symbol)
		err := where.First(&m).Error
		if err == gorm.ErrRecordNotFound {
			global.DB.Model(CombineOrder{}).Save(&v)
		}
		// else if err == nil {
		// 	where.Update("commission", v.Commission) // 每次修改字段的时候修改
		// }
	}
}

func UpsertAccountTrade(data AccountTrade) {
	m := AccountTrade{}
	where := global.DB.Model(AccountTrade{}).Where("id", data.ID)
	err := where.First(&m).Error
	if err == gorm.ErrRecordNotFound {
		global.DB.Model(AccountTrade{}).Create(data)
	}
}

func UpsertOrder(data Order) {
	m := Order{}
	where := global.DB.Model(Order{}).Where("orderId", data.OrderId)
	err := where.First(&m).Error
	if err == gorm.ErrRecordNotFound {
		global.DB.Model(Order{}).Create(data)
	}
}

func GetAllPlan() ([]Plan, error) {
	list := make([]Plan, 0)
	err := global.DB.Model(Plan{}).Find(&list).Error
	return list, err
}

func GetPlanByUserId(userId int) ([]Plan, error) {
	list := make([]Plan, 0)
	err := global.DB.Model(Plan{}).Where("userId", userId).Find(&list).Error
	return list, err
}

func GetPlanById(id int64) (Plan, error) {
	data := Plan{}
	err := global.DB.Model(Plan{}).Where("id", id).First(&data).Error
	return data, err
}

func GetUserPlanById(userid int, id int64) (Plan, error) {
	data := Plan{}
	err := global.DB.Model(Plan{}).Where("userId", userid).Where("id", id).First(&data).Error
	return data, err
}

func SavePlan(data Plan) error {
	return global.DB.Model(data).Save(data).Error
}

func UpdatePlan(userId int, data Plan) error {
	return global.DB.Model(data).Where("userId", userId).Save(data).Error
}

func CreatePlan(data *Plan) error {
	return global.DB.Model(data).Create(&data).Error
}

func DelPlan(id int64) error {
	return global.DB.Where("id", id).Delete(&Plan{}).Error
}

func DonePlan(id int64) error {
	return global.DB.Model(Plan{}).Where("id", id).Update("status", 1).Error
}

func GetAllUser() ([]User, error) {
	list := make([]User, 0)
	err := global.DB.Model(User{}).Find(&list).Error
	return list, err
}

func GetUserByToken(token string) (User, error) {
	user := User{}
	err := global.DB.Model(User{}).Where("token", token).First(&user).Error
	return user, err
}

func GetUserByUserId(userId int) (User, error) {
	user := User{}
	err := global.DB.Model(User{}).Where("id", userId).First(&user).Error
	return user, err
}
