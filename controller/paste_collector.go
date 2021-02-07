package controller

import "github.com/gin-gonic/gin"

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

}
