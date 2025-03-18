package main

import (
	"fmt"
	Config "helptrade/config"
	"helptrade/controller"
	"helptrade/dao"
	"helptrade/global"
	"helptrade/service"
	"net/http"
	"time"

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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,token")
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}

		userInfo, err := dao.GetUserByToken(token) // todo 改成redis或者jwttoken
		fmt.Println("userInfo", userInfo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}
		c.Set("user", userInfo)

		// 继续处理后续请求
		c.Next()
	}
}
func main() {

	go func() {
		for {
			service.FetchAndCombineAccountTrade()
			time.Sleep(30 * time.Minute)
		}
	}()
	// go func() {
	// 	for {
	// 		log.Println("doPlan")
	// 		service.DoPlan()
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	r := gin.Default()

	r.Static("/static", "./public")

	api := r.Group("/api")
	api.Use(CORSMiddleware()).Use(AuthMiddleware())
	{
		api.GET("/getCombineOrderList", controller.GetCombineOrderList)
		api.POST("/editCommnet", controller.EditCommnet)

		api.GET("/getPlanList", controller.GetPlanList)
		api.POST("/savePlan", controller.SavePlan)
		api.POST("/delPlan", controller.DelPlan)
	}

	r.Run()
}
