package riaktimeline

import (
	"encoding/json"
	"io/ioutil"
)

type TimelineConfig struct {
	Nodes              []string `json:"nodes"`
	LogBucket          string   `json:"log-bucket"`
	IndexBucket        string   `json:"index-bucket"`
	DaysToKeep         int      `json:"days-to-keep"`
	EnableDebugLogging bool     `json:"enable-debug"`
}

var defaultConfig = TimelineConfig{
	Nodes:              []string{"127.0.0.1:11087"},
	LogBucket:          "log-entries",
	IndexBucket:        "log-indexes",
	DaysToKeep:         30,
	EnableDebugLogging: true,
}

func ReadConfig(path string) (TimelineConfig, error) {
	var config TimelineConfig
	file, e := ioutil.ReadFile(path)
	if e != nil {
		return config, e
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return config, err
	}
	return config, nil
}
