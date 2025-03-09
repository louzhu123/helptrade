package main

import (
	"fmt"
	Config "helptrade/config"
	"helptrade/controller"
	"helptrade/global"
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
	// service.FetchAndSaveAllAccountTrade()
	// list := service.CombineAccountTrade()
	// dao.SaveCombineOrder(list)
	// return
	// go func() {
	// 	for {
	// 		service.FetchAndCombineOrder()
	// 		time.Sleep(15 * time.Minute)
	// 	}
	// }()

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/getCombineOrderList", controller.GetCombineOrderList)
	r.POST("/editCommnet", controller.EditCommnet)

	r.Run()
}
