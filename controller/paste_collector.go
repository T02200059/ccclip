package controller

import (
	"ccclip/libs"
	h "ccclip/pkg/restful"
	"github.com/atotto/clipboard"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PasteCollectorController struct {
}

func NewPasteCollectorControllerProvider() *PasteCollectorController {
	return &PasteCollectorController{}
}

func (pcc *PasteCollectorController) PingPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (pcc *PasteCollectorController) HandlePaste(c *gin.Context) {
	var req struct {
		Payload string `json:"payload"`
	}

	var err error

	err = c.ShouldBind(&req)
	if err != nil {
		log.Error(err)
		h.SendError(c, err, nil)
		return
	}

	err = handlePaste(req.Payload, "")
	if err != nil {
		log.Error(err)
		h.SendError(c, err, nil)
		return
	}

	h.SendOK(c, nil)
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

	err = clipboard.WriteAll(tcc)
	if err != nil {
		return
	}

	return
}
