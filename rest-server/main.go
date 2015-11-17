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

func init() {
	flag.StringVar(&configPath, "c", "", "path to config json")
}

func main() {
	flag.Parse()
	rtl.Bootstrap(configPath)
	router := gin.Default()
	router.GET("/v1/log/:name", func(c *gin.Context) {
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

	address := fmt.Sprintf(":%d", rtl.Config.RestPort)
	router.Run(address) // listen and serve on 0.0.0.0:8080
}
