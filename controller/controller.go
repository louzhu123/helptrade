package controller

import (
	"fmt"
	"helptrade/dao"
	"helptrade/global"
	"helptrade/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCombineOrderList(c *gin.Context) {

	var req global.GetCombineOrderListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)

	list, _ := dao.QueryCombineOrder(req)

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

func GetPlanList(c *gin.Context) {
	// var req global.GetPlanListReq
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }

	list, _ := service.GetPlanList()

	c.JSON(http.StatusOK, gin.H{
		"message": list,
	})
}

func SavePlan(c *gin.Context) {
	var req global.SavePlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.SavePlan(req)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func DelPlan(c *gin.Context) {
	var req global.DelPlanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	service.DelPlan(req)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
