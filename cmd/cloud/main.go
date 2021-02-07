package main

import (
	"ccclip/controller"
	"ccclip/libs"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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
	log.Info("Cloud collector.")

	var err error
	var errs = make(chan error, 3)
	var m = make(map[libs.UserCode]*libs.ClipRecord)

	pcc := controller.NewPasteCollectorControllerProvider(m)

	r := gin.Default()
	r.GET("/ping", pcc.PingPong)
	r.POST("/async", pcc.HandlePaste)
	log.Info("Paste collector server start.")
	go func() {
		errs <- r.Run(libs.CloudPort)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	select {
	case <-c:
		err = errors.New("input ^C")
	case err = <-errs:
	}

	log.Errorf("exited by: %+v", err)

	log.Info("Bye bye.")
}
