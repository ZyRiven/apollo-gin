package main

import (
	"apollo/httpServer"
	"apollo/setting"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	httpserver.RouterWeb(setting.HttpPort)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signalChan:
		setting.ZAPS.Infof("程序退出")
		setting.CronProcess.Stop()
		return
	}
}
