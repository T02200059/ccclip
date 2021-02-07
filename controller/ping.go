package controller

import (
	"ccclip/libs"
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type PingController struct {
}

func NewPingControllerProvider() *PingController {
	return &PingController{}
}

func (pc *PingController) Run(ctx context.Context) (err error) {
	log.Info("Ping pong service start.")
	var errs = make(chan error, 2)

	for {
		err = doPingPong()
		if err != nil {
			errs <- err
			break
		}

		time.Sleep(libs.PingPongInterval)
	}

	select {
	case <-ctx.Done():
	case <-errs:
	}

	log.Info("Ping pong service stopped.")
	return
}

// 向 cloud 发送 ping 信息, 表明本机处于开机状态, 需要同步剪贴板. 同时
func doPingPong() (err error) {
	// TODO
	return
}
