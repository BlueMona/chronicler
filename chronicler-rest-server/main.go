package main

import (
	"flag"
	"fmt"
	entity "github.com/PeerioTechnologies/chronicler/entity"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
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
	bootstrap(configPath)
	router := gin.Default()

	router.GET("/v1/log/get/:name", func(c *gin.Context) {
		name := c.Param("name")
		var msg struct {
			Error   error                `json:"error"`
			Payload entity.TimelineIndex `json:"payload"`
		}
		result, err := dao.GetTimeline(name)
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
			var msg struct {
				Error error `json:"error"`
			}
			status := http.StatusOK
			if err := dao.SaveLog(frm.Id, frm.Level, frm.Type, frm.Log); err != nil {
				status = http.StatusBadRequest
				msg.Error = err
			}
			c.JSON(status, msg)
		}
	})

	router.GET("/v1/report", func(c *gin.Context) {
		var msg struct {
			State       string `json:"state"`
			Description string `json:"description"`
		}
		state, desc := dao.ReportState()
		msg.State = state
		msg.Description = desc
		c.JSON(http.StatusOK, msg)
	})

	address := fmt.Sprintf(":%d", config.RestPort)
	router.Run(address)

	//Handling iterrupt
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...\n")
			// cleanup(services, c)
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
