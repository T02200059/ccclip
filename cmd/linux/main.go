package main

import (
	"ccclip/controller"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	fmt.Println("Use package: ", "github.com/atotto/clipboard")

	var err error
	var errs = make(chan error, 3)
	ctx, cancel := context.WithCancel(context.Background())

	// sqlite3

	// copy
	ccc := controller.NewCopyCollectorControllerProvider()
	go func() {
		errs <- ccc.Run(ctx)
	}()

	// paste
	pcc := controller.NewPasteCollectorControllerProvider()

	r := gin.Default()
	r.GET("/ping", pcc.PingPong)
	r.POST("/paste", pcc.HandlePaste)

	log.Info("Paste collector server start.")
	go func() {
		errs <- r.Run(":22122")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	select {
	case <-c:
		err = errors.New("input ^C")
	case err = <-errs:
	}

	log.Errorf("exited by: %+v", err)

	cancel() // close ctx

	log.Info("Wait 1 second...")
	time.Sleep(1 * time.Second)

	// TODO: CLOSE GIN

	log.Info("Bye bye.")
}
