package httpserver

import (
	"apollo/httpServer/controller"
	"apollo/report"
	"apollo/setting"
	"time"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)


func init() {
	// 初始化时间：东8
	var cstZone = time.FixedZone("CST", 8*3600)
	time.Local = cstZone
	fmt.Println("                    _ _       \n   __ _ _ __   ___ | | | ___  \n  / _` | '_ \\ / _ \\| | |/ _ \\ \n | (_| | |_) | (_) | | | (_) |\n  \\__,_| .__/ \\___/|_|_|\\___/ \n       |_|                    ")
	setting.GetConf()
	setting.InitLogger()
	setting.StreamInit()
	report.ReportServiceInit()
}

func RouterWeb(port string) {
	r := gin.New()
	r.Use(gin.Recovery())
	gin.SetMode(setting.AppMode)

	loginRouter := r.Group("/system")
	{
		loginRouter.GET("/ping", controller.Ping)
	}
	sendR := r.Group("/send")
	{
		// sendR.POST("/sensor")
		// sendR.POST("/action")
		sendR.POST("/modbus",controller.SendModbus)
	}

	setting.ZAPS.Infof("gin 初始化端口 %s!", port)
	go func() {
		err := r.Run(port)
		if err != nil {
			setting.ZAPS.Errorf("gin启动错误 %v", err)
			os.Exit(0)
		}
	}()
}
