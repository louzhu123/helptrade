package main

import (
	"fmt"
	Config "main/config"
	"main/controller"
	"main/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	global.Cfg = &Config.Config{}
	err := ini.MapTo(global.Cfg, "config.ini")
	if err != nil {
		panic(err)
	}

	fmt.Println(global.Cfg.ConfigMysql.Dsn)
	initMysql(global.Cfg.ConfigMysql.Dsn)
}

func initMysql(dsn string) {
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open mysql failed,", err)
	}
	global.DB = d
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许的请求来源
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// 设置允许的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// 设置允许的请求凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 如果是 OPTIONS 请求，直接返回
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// 继续处理其他请求
		c.Next()
	}
}

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	r.Use(CORSMiddleware())

	// 定义一个 GET 接口
	r.GET("/getCombineOrderList", controller.GetCombineOrderList)

	// 启动服务器，默认在0.0.0.0:8080启动服务
	r.Run()
	// dao.SelectAccountTradeGroupByOrderId()
	// list := dao.CombineAccountOrder()
	// dao.SaveCombineOrder(list)
	// futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)
	// res, err := futuresClient.newtrade().Do(context.Background())
	// fmt.Println(err)
	// fmt.Println(res)
	// res, err := futuresClient.NewListOrdersService().Do(context.Background())

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(res)
	// for _, v := range res {
	// 	//fmt.Printf("%+v", v)
	// 	//j, _ := json.Marshal(v)
	// 	saveErr := DB.Model(Order{}).Create(&v).Error
	// 	if saveErr != nil {
	// 		fmt.Println(err)
	// 	}
	// }

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

	// 查询根据orderid group
}
