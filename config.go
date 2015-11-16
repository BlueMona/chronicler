package riaktimeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TimelineConfig struct {
	Nodes       []string `json:"nodes"`
	LogBucket   string   `json:"log-bucket"`
	IndexBucket string   `json:"index-bucket"`
	DaysToKeep  int      `json:"days-to-keep"`
}

var defaultConfig = TimelineConfig{
	Nodes:       []string{"127.0.0.1:11087"},
	LogBucket:   "log-entries",
	IndexBucket: "log-indexes",
	DaysToKeep:  30,
}

func ReadConfig(path string) TimelineConfig {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("Config reading error:", e)
		fmt.Println("Using defaut configuration...")
		return defaultConfig
	}
	var config TimelineConfig
	json.Unmarshal(file, &config)
	return config
}
