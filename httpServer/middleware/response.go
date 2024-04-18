package middleware

import "github.com/gin-gonic/gin"

type Gin struct {
	C *gin.Context
}

func (g *Gin) Success(code int, data interface{}, msg ...string) {
	g.C.JSON(200, gin.H{
		"Code":    code,
		"Message": msg,
		"Data":    data,
	})
}

func (g *Gin) Error(code int, data interface{}, msg ...string) {
	g.C.JSON(400, gin.H{
		"Code":    code,
		"Message": msg,
		"Data":    data,
	})
}