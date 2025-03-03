package controller

import (
	"context"
	"fmt"
	"main/dao"
	"main/global"
	"net/http"

	"github.com/adshao/go-binance/v2"
	"github.com/gin-gonic/gin"
)

// 懒得写servcie层了
func GetCombineOrderList(c *gin.Context) {

	list, _ := dao.QueryCombineOrder()

	c.JSON(http.StatusOK, gin.H{
		"message": list,
	})
}

func EditCommnet(c *gin.Context) {
	var req global.EditCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	dao.UpdateCombineOrderComment(req.Id, req.Comment)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func FetchTradeData() {
	futuresClient := binance.NewFuturesClient(global.Cfg.ApiKey, global.Cfg.SecretKey)

	res, err := futuresClient.NewListOrdersService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res {
		saveErr := global.DB.Model(dao.Order{}).Create(&v).Error
		if saveErr != nil {
			fmt.Println(err)
		}
	}

	res1, err := futuresClient.NewListAccountTradeService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range res1 {
		saveErr := global.DB.Model(dao.AccountTrade{}).Create(&v).Error
		if saveErr != nil {
			fmt.Println(err)
			break
		}
	}

	list := dao.CombineAccountOrder()
	dao.SaveCombineOrder(list)
}
