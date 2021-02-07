package controller

import (
	"ccclip/libs"
	"ccclip/pkg/restful"
	"context"
	"encoding/json"
	"errors"
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

func (ccc *CopyCollectorController) UserCode() (result libs.UserCode) {
	return ccc.BasicUserCode + ccc.Platform
}

// Run: Block here. Use 'go run' usually.
func (ccc *CopyCollectorController) Run(ctx context.Context) (err error) {
	log.Info("Clipboard copycat start.")
	var errs = make(chan error, 2)

	var last string
	var current string
	var curRecord = &libs.ClipRecord{
		User:      ccc.UserCode(),
		UpdatedAt: time.Now(),
	}

	for {
		// TODO: 将读写剪贴板的方法与操作系统隔离
		// read clipboard
		current, err = clipboard.ReadAll()
		if err != nil {
			log.Error(err)
			errs <- err
			break
		}

		// trim '\n'
		tcc := libs.DefaultTrimmer(current)

		// compare
		if last != tcc {
			log.Debugf("A new local clipboard payload: %s", tcc)
			curRecord.Payload = tcc
			curRecord.UpdatedAt = time.Now()
		}
		last = tcc

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
			log.Error(err)
			errs <- err
			break
		}

		// next...
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
func (ccc *CopyCollectorController) queryClipboard(record *libs.ClipRecord) (response string, err error) {
	if record == nil {
		err = errors.New("record cannot be nil")
		return
	}

	var postUrl string

	postUrl = ccc.CloudURL + libs.SuffixAsync

	resp, err := restful.DoPost(postUrl, record, nil, nil)
	if err != nil {
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	respCopy := &restful.CopyResponse{}
	err = json.Unmarshal(respBody, respCopy)
	if err != nil {
		return
	}
	response = respCopy.Data
	return
}

// 对 cloud 同步来的剪贴板内容进行检查并写入到本地剪贴板.
func handlePaste(payload string, current string) (err error) {
	if payload == "" {
		return // nothing happened
	}

	tcc := libs.DefaultTrimmer(payload)
	if tcc == current {
		return // nothing happened.
	}

	log.WithField("payload", payload).Info("Update clipboard on this device.")
	err = clipboard.WriteAll(tcc)
	if err != nil {
		return
	}

	return
}
