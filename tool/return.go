package tool

import "github.com/gin-gonic/gin"

func JsonSuccess(c *gin.Context,data interface{})  {
	c.JSON(200,gin.H{
		"code":200,
		"msg":"成功",
		"data":data,
	})
}

func JsonError(c *gin.Context,msg string)  {
	c.JSON(200,gin.H{
		"code":400,
		"msg":msg,
	})
}