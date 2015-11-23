package riaktimeline

import (
	"encoding/json"
	"io/ioutil"
)

type TimelineDAOConfig struct {
	Nodes              []string `json:"nodes"`
	LogBucket          string   `json:"log-bucket"`
	IndexBucket        string   `json:"index-bucket"`
	DaysToKeep         int      `json:"days-to-keep"`
	EnableDebugLogging bool     `json:"enable-debug"`
}

var defaultConfig = TimelineDAOConfig{
	Nodes:              []string{"127.0.0.1:11087"},
	LogBucket:          "log-entries",
	IndexBucket:        "log-indexes",
	DaysToKeep:         30,
	EnableDebugLogging: true,
}
