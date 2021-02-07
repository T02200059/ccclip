package main

import (
	"ccclip/controller"
	"context"
	"errors"
	"fmt"
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

const cloudURL = "http://localhost:22123"

func main() {
	log.Info("Hello world.")
	log.Info("Linux clipboard collector.")
	fmt.Println("Use package: ", "github.com/atotto/clipboard")

	var err error
	var errs = make(chan error, 3)
	ctx, cancel := context.WithCancel(context.Background())

	// copy & paste 在无公网 ip 的环境下, 通过主动轮询来刷新数据.
	ccc := controller.NewCopyCollectorControllerProvider(cloudURL)
	go func() {
		errs <- ccc.Run(ctx)
	}()

	// 在非 cloud 中, 因为当前设备均无公网 ip, 故不使用主动服务器
	// pcc := controller.NewPasteCollectorControllerProvider()
	//
	// r := gin.Default()
	// r.GET("/ping", pcc.PingPong)
	// r.POST("/paste", pcc.HandlePaste)
	// log.Info("Paste collector server start.")
	// go func() {
	// 	errs <- r.Run(":22122")
	// }()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	select {
	case <-c:
		err = errors.New("input ^C")
	case err = <-errs:
	}

	log.Errorf("exited by: %+v", err)

	cancel() // close ctx

	log.Info("Bye bye.")
}
