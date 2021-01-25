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
	TestStucker.Purchase()
	//if TestStucker.Counter() > 0 {
	//TestStucker.Purchase()
	//}
}

var TestStucker = NewStuck(3000)

func main() {
	r := setUpRouter()

	r.Run(":9999")
}
