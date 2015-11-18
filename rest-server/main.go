package main

import (
	"flag"
	"fmt"
	rtl "github.com/PeerioTechnologies/riak-timeline-service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	configPath string
)

type StoreLogRequest struct {
	Id    string `form:"id" json:"id"  binding:"required"`
	Level string `form:"level" json:"level"  binding:"required"`
	Type  string `form:"type" json:"type"  binding:"required"`
	Log   string `form:"log" json:"log"  binding:"required"`
}

func init() {
	flag.StringVar(&configPath, "c", "", "path to config json")
}

//TODO - Errors granularity.
//TODO - separate routing and handlers
//
func main() {
	flag.Parse()
	rtl.Bootstrap(configPath)
	router := gin.Default()

	router.GET("/v1/log/get/:name", func(c *gin.Context) {
		name := c.Param("name")
		var msg struct {
			Error   error             `json:"error"`
			Payload rtl.TimelineIndex `json:"payload"`
		}
		result, err := rtl.FetchMergedTimeline(name)
		status := http.StatusOK
		if err != nil {
			status = http.StatusInternalServerError
		}
		msg.Error = err
		msg.Payload = result
		c.JSON(status, msg)
	})

	router.POST("/v1/log/append", func(c *gin.Context) {
		var frm StoreLogRequest
		// This will infer what binder to use depending on the content-type header.
		if c.Bind(&frm) == nil {
			//TODO validate input
			status := http.StatusOK
			err := rtl.SaveLog(frm.Id, frm.Level, frm.Type, frm.Log)
			if err != nil {
				status = http.StatusBadRequest
			}
			var msg struct {
				Error error `json:"error"`
			}
			msg.Error = err
			c.JSON(status, msg)
		}
	})

	address := fmt.Sprintf(":%d", rtl.Config.RestPort)
	router.Run(address)
}
