package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
}

func main() {
	log.Info("Hello world.")
	log.Info("Linux clipboard collector.")

	var err error
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/")
	log.Info("Collector start.")
	err = r.Run(":22122")

	log.Errorf("Collector Stopped. err=%+v", err)
}
