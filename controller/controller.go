package controller

import (
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

	userInfo, _ := c.Get("user")
	user := userInfo.(dao.User)
	list, _ := dao.QueryCombineOrder(user.Id, req)

	c.JSON(http.StatusOK, gin.H{
		"message": list,
	})
}

func GetCombineOrderStatis(c *gin.Context) {

	var req global.GetCombineOrderListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userInfo, _ := c.Get("user")
	user := userInfo.(dao.User)
	list, _ := dao.GetCombineOrderStatis(user.Id, req)

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

	userInfo, _ := c.Get("user")
	user := userInfo.(dao.User)

	dao.UpdateCombineOrderComment(int64(user.Id), req.Id, req.Comment)

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

	userInfo, _ := c.Get("user")
	user := userInfo.(dao.User)

	list, _ := service.GetPlanList(user.Id)

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

	userInfo, _ := c.Get("user")
	user := userInfo.(dao.User)
	service.SavePlan(user.Id, req)

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
