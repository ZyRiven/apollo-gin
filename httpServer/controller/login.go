package controller

import (
	"apollo/httpServer/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	appG := middleware.Gin{C: c}
	t := time.Now()
	data := map[string]interface{}{
		"data": "ping",
		"time": t,
		"timeUnix": t.Unix(),
	}
	appG.Success(200, data)
}
