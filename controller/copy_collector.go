package controller

import (
	"ccclip/libs"
	"ccclip/pkg/restful"
	"context"
	"encoding/json"
	"github.com/atotto/clipboard"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type CopyCollectorController struct {
	CloudURL      string // base url
	BasicUserCode libs.UserCode
	Platform      libs.UserCode
}

func NewCopyCollectorControllerProvider(cloudURL string, b, p libs.UserCode) *CopyCollectorController {
	return &CopyCollectorController{
		CloudURL:      cloudURL,
		BasicUserCode: b,
		Platform:      p,
	}
}

type ClipRecord struct {
	User      libs.UserCode `json:"user"`
	Payload   string        `json:"payload"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (ccc *CopyCollectorController) UserCode() (result libs.UserCode) {
	return ccc.BasicUserCode + ccc.Platform
}

// Run: Block here. Use 'go run' usually.
func (ccc *CopyCollectorController) Run(ctx context.Context) (err error) {
	log.Info("Clipboard copycat start.")
	var errs = make(chan error, 2)

	var last string
	var current string

	for {
		var curRecord *ClipRecord

		// TODO: 将读写剪贴板的方法与操作系统隔离
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
				User:      ccc.UserCode(),
				Payload:   tcc,
				UpdatedAt: time.Now(),
			}
		}

		// async
		var latestCopy string
		latestCopy, err = ccc.queryClipboard(curRecord)
		if err != nil {
			log.Error(err)
			time.Sleep(libs.NetworkErrInterval) // 这部分错误不强制退出.
		}

		// paste
		err = handlePaste(latestCopy, current)
		if err != nil {
			errs <- err
			break
		}

		// next...
		last = tcc
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
func (ccc *CopyCollectorController) queryClipboard(record *ClipRecord) (response string, err error) {
	if record != nil {
		log.Infof("Async a new clipboard record: %+v", record)
	}

	var postUrl string
	var postBody []byte

	postUrl = ccc.CloudURL + libs.SuffixAsync

	if record != nil {
		postBody, err = json.Marshal(record)
		if err != nil {
			log.Error(err)
			return
		}
	}

	resp, err := restful.DoPost(postUrl, postBody, nil, nil)
	if err != nil {
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	response = string(respBody)
	return
}
