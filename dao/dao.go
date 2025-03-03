package dao

import (
	"helptrade/global"

	"github.com/adshao/go-binance/v2/futures"
	"gorm.io/gorm"
)

func UpdateCombineOrderComment(id int64, comment string) {
	global.DB.Table("combine_order").Where("id", id).Update("comment", comment)
}

func QueryCombineOrder() ([]CombineOrder, error) {
	var list []CombineOrder
	err := global.DB.Model(&CombineOrder{}).Order("startTime desc").Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
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

func GetAllOrder() ([]Order, error) {
	list := make([]Order, 0)
	err := global.DB.Model(Order{}).Find(&list).Order("time asc").Error
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

func UpsertAccountTrade(data *futures.AccountTrade) {
	m := AccountTrade{}
	where := global.DB.Model(AccountTrade{}).Where("id", data.ID)
	err := where.First(&m).Error
	if err == gorm.ErrRecordNotFound {
		global.DB.Model(AccountTrade{}).Save(data)
	}
}

func UpsertOrder(data *futures.Order) {
	m := Order{}
	where := global.DB.Model(Order{}).Where("orderId", data.OrderID)
	err := where.First(&m).Error
	if err == gorm.ErrRecordNotFound {
		global.DB.Model(Order{}).Save(data)
	}
}
