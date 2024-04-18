package controller

import (
	"apollo/httpServer/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SendModbus(c *gin.Context) {
	appG := middleware.Gin{C: c}
	var data map[string]interface{}
	
	if err := c.BindJSON(&data); err != nil {
		appG.Error(400, err)
		return
	}
	fmt.Println(data)
	appG.Success(200, "ping")
}
