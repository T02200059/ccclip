package main

import (
	"ccclip/controller"
	"ccclip/libs"
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetReportCaller(true)
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
	defer cancel()

	// copy & paste 在无公网 ip 的环境下, 通过主动轮询来刷新数据.
	ccc := controller.NewCopyCollectorControllerProvider(libs.CloudURL, libs.OriginUserYTB, libs.GetPlatformCode())
	go func() {
		errs <- ccc.Run(ctx)
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
