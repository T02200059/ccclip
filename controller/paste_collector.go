package controller

import (
	"ccclip/libs"
	h "ccclip/pkg/restful"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PasteCollectorController struct {
	sessions map[libs.UserCode]*libs.ClipRecord
}

func NewPasteCollectorControllerProvider(m map[libs.UserCode]*libs.ClipRecord) *PasteCollectorController {
	return &PasteCollectorController{
		sessions: m,
	}
}

func (pcc *PasteCollectorController) PingPong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (pcc *PasteCollectorController) HandlePaste(c *gin.Context) {
	var req libs.ClipRecord
	var err error
	err = c.ShouldBind(&req)
	if err != nil {
		log.Error(err)
		h.SendError(c, err, nil)
		return
	}

	fmt.Printf("%+v\n", req)

	cp := pcc.async(&req)

	h.SendOK(c, cp) // TODO: COMPARE THE LATEST COPY
}

func (pcc *PasteCollectorController) async(info *libs.ClipRecord) (latestCopy string) {
	if info == nil {
		return
	}

	if pcc.sessions == nil {
		log.Panic("Sessions map is nil!")
	}

	latestCopy = info.Payload

	u := libs.DecodeUser(info.User)
	if s, ok := pcc.sessions[u]; ok {
		if s.UpdatedAt.After(info.UpdatedAt) {
			latestCopy = s.Payload
			return
		}
	}

	pcc.sessions[u] = info

	return
}
