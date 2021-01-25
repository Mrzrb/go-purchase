package main

import (
	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/purchase", Purchase)
	return r
}

func Purchase(c *gin.Context) {
	if TestStucker.Purchase() {
		c.JSON(200, "购买成功")
	} else {
		c.JSON(200, "购买失败")
	}
}

var TestStucker = NewStuck(3000)

func main() {
	r := setUpRouter()

	r.Run(":9999")
}
