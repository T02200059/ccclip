package restful

import "github.com/gin-gonic/gin"

func SendOK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": "Success",
		"data": data,
		"msg":  "",
	})
}

func SendError(c *gin.Context, err error, data interface{}) {
	c.JSON(200, gin.H{
		"code": "Error",
		"data": data,
		"msg":  err.Error(),
	})
}
