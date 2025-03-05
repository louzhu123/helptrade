package controller

import (
	"fmt"
	"helptrade/dao"
	"helptrade/global"
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
