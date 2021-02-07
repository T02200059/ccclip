package main

import (
	"ccclip/controller"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	log.Info("Hello world.")
	log.Info("Linux clipboard collector.")

	// sqlite3

	// copy
	ccc := controller.NewCopyCollectorControllerProvider()
	go ccc.Run()

	// paste
	var err error

	pcc := controller.NewPasteCollectorControllerProvider()

	r := gin.Default()
	r.GET("/ping", pcc.PingPong)
	r.POST("/paste", pcc.HandlePaste)

	log.Info("Paste collector server start.")
	err = r.Run(":22122")
	if err != nil {
		log.Errorf("Collector Stopped. err=%+v", err)
	}

	log.Info("Bye bye.")
}
