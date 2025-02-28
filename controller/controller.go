package controller

import (
	"main/dao"
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
