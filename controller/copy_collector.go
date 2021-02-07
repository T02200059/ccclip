package controller

import (
	"ccclip/libs"
	"context"
	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"time"
)

type CopyCollectorController struct {
	CloudURL string // base url
}

func NewCopyCollectorControllerProvider(cloudURL string) *CopyCollectorController {
	return &CopyCollectorController{
		CloudURL: cloudURL,
	}
}

type ClipRecord struct {
	Payload   string    `json:"payload"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Run: Block here. Use 'go run' usually.
func (ccc *CopyCollectorController) Run(ctx context.Context) (err error) {
	log.Info("Clipboard copycat start.")
	var errs = make(chan error, 2)

	var last string
	var current string

	for {
		var curRecord *ClipRecord

		// read clipboard
		current, err = clipboard.ReadAll()
		if err != nil {
			errs <- err
			break
		}

		// trim '\n'
		tcc := libs.DefaultTrimmer(current)

		// compare
		if last != tcc {
			curRecord = &ClipRecord{
				Payload:   tcc,
				UpdatedAt: time.Now(),
			}
		}

		// async
		var latestCopy string
		latestCopy, err = queryClipboard(curRecord)
		if err != nil {
			log.Error(err)
			continue
		}

		// paste
		err = handlePaste(latestCopy, current)
		if err != nil {
			log.Error(err)
			continue
		}

		// next...
		last = current
		time.Sleep(libs.CopyCollectorInterval)
	}

	select {
	case <-ctx.Done():
	case <-errs:
	}

	log.Info("Clipboard copycat stopped.")
	return
}

// POST 请求发送剪贴板记录并获取最新同步.
func queryClipboard(record *ClipRecord) (response string, err error) {
	if record != nil {
		log.Infof("Async a new clipboard record: %+v", record)
	}

	// TODO: request

	return
}
