package riaktimeline

import (
	"fmt"
	riak "github.com/basho/riak-go-client"
	gosnow "github.com/sdming/gosnow"
	"strconv"
	"sync"
	"time"
)

func Bootstrap(configPath string) {
	logInfo("[Bootstrap]", "Loading config from path \"%s\"", configPath)
	//reading configuration
	if config, err := ReadConfig(configPath); err == nil {
		Config = config
	} else {
		logError("[Bootstrap]", "Error reading config \"%s\", using default", configPath)
		Config = defaultConfig
	}
	//connect to Riak
	nodeOptions := buildNodeOptions(Config.Nodes)
	RiakCluster = initCluster(nodeOptions)
	if err := RiakCluster.Start(); err != nil {
		logErr("Error connection to Riak cluster %v", err)
		panic(fmt.Sprintf("Error %v", err))
	}
	initIndexBucket(Config.IndexBucket)
	//init SnowFlakes
	snowFlake, _ = gosnow.Default()
	//enable debug
	EnableDebugLogging = Config.EnableDebugLogging
	riak.EnableDebugLogging = Config.EnableDebugLogging
}

func FetchMergedTimeline(id string) (TimelineIndex, error) {
	var index TimelineIndex
	var err error
	index, err = getTimeline(id)
	if err != nil {
		return nil, err
	}
	//TODO add removal of old entries by channel
	index, _ = splitByDaysAge(index, Config.DaysToKeep)
	//async fill IndexEntry.Caption with real log data
	itemsQty := len(index)
	group := new(sync.WaitGroup)
	group.Add(itemsQty)
	for i := 0; i < itemsQty; i++ {
		go func(entry *IndexEntry, group *sync.WaitGroup) {
			storedMsg, _ := fetchLog(entry.Key)
			entry.Caption = storedMsg
			group.Done()
		}(&index[i], group)
	}
	group.Wait()
	return index, nil
}

func SaveLog(userId string, level string, typeStr string, msg string) error {
	id, _ := snowFlake.Next()
	idStr := strconv.FormatUint(id, 10)
	entry := IndexEntry{
		Key:     idStr,
		Time:    time.Now(),
		Level:   level,
		Type:    typeStr,
		Caption: "",
	}
	errc := make(chan error)
	go func() {
		errc <- storeLog(idStr, msg)
	}()
	go func() {
		errc <- appendToTimeline(userId, entry)
	}()
	var err error
	for i := 0; i < 2; i++ {
		if e := <-errc; e != nil {
			err = e
		}
	}
	return err
}
