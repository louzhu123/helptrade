package main

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/go-ini/ini"
)

func init() {
	Cfg = &Config{}
	err := ini.MapTo(Cfg, "config.ini")
	if err != nil {
		panic(err)
	}

	fmt.Println(Cfg.ConfigMysql.Dsn)
	initMysql(Cfg.ConfigMysql.Dsn)
}

func main() {
	futuresClient := binance.NewFuturesClient(Cfg.ApiKey, Cfg.SecretKey)
	res, err := futuresClient.NewListOrdersService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	for _, v := range res {
		//fmt.Printf("%+v", v)
		//j, _ := json.Marshal(v)
		saveErr := DB.Model(Order{}).Create(&v).Error
		if saveErr != nil {
			fmt.Println(err)
		}
	}

	// {"buyer":true,"commission":"0.00301236","commissionAsset":"USDT","id":32684952,"maker":false,"orderId":221465321,"price":"1.544800","qty":"3.9","quoteQty":"6.0247200","realizedPnl":"-0.03595577","side":"BUY","positionSide":"SHORT","symbol":"KAITOUSDT","time":1740295944617}
	// 	res, err := futuresClient.NewListAccountTradeService().Do(context.Background())
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	// fmt.Println(res)
	// 	for _, v := range res {
	// 		//fmt.Printf("%+v", v)
	// 		j, _ := json.Marshal(v)
	// 		fmt.Println(string(j))
	// 		saveErr := DB.Model(AccountTrade{}).Create(&v).Error
	// 		if saveErr != nil {
	// 			fmt.Println(err)
	// 			break
	// 		}
	// 	}
}
