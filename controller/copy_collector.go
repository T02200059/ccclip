package controller

import (
	"ccclip/libs"
	"context"
	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"time"
)

type CopyCollectorController struct {
}

func NewCopyCollectorControllerProvider() *CopyCollectorController {
	return &CopyCollectorController{}
}

// Run: Block here. Use 'go run' usually.
func (ccc *CopyCollectorController) Run(ctx context.Context) (err error) {
	log.Info("Clipboard copycat start.")
	var errs = make(chan error, 2)

	var last string
	var current string

	for {
		// 读取剪切板中的内容到字符串, 本身不会读取空值.
		current, err = clipboard.ReadAll()
		if err != nil {
			errs <- err
			break
		}

		if last != current {
			// send to cloud server
			ok := handleCurrentCopy(current)
			if ok {
				last = current
			}
		}

		time.Sleep(libs.CopyCollectorInterval)
	}

	select {
	case <-ctx.Done():
	case <-errs:
	}

	log.Info("Clipboard copycat stopped.")
	return
}

// 负责整理字符串并发送到同步 cloud.
func handleCurrentCopy(cc string) (ok bool) {
	tcc := libs.DefaultTrimmer(cc)
	println(tcc)

	// trim '\n'

	// TODO: SEND TO CLOUD SERVER

	ok = true
	return
}
