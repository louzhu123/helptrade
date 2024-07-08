package main

import (
	"context"
	"encoding/json"
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
	//res, err := futuresClient.NewListOrdersService().Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)
	//for _, v := range res {
	//	//fmt.Printf("%+v", v)
	//	//j, _ := json.Marshal(v)
	//	saveErr := DB.Model(TradeRecord{}).Create(&v).Error
	//	if saveErr != nil {
	//		fmt.Println(err)
	//	}
	//}

	res, err := futuresClient.NewListAccountTradeService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	for _, v := range res {
		//fmt.Printf("%+v", v)
		j, _ := json.Marshal(v)
		fmt.Println(string(j))
		//saveErr := DB.Model(TradeRecord{}).Create(&v).Error
		//if saveErr != nil {
		//	fmt.Println(err)
		//}
	}
}
