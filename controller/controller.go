package controller

import (
	"main/dao"
	"main/global"
	"net/http"

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
