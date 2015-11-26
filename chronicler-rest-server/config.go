package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServiceConfig struct {
	Nodes              []string `json:"nodes"`
	LogBucket          string   `json:"log-bucket"`
	IndexBucket        string   `json:"index-bucket"`
	DaysToKeep         int      `json:"days-to-keep"`
	EnableDebugLogging bool     `json:"enable-debug"`
	RestPort           int      `json:"rest-port"`
}

var defaultConfig = ServiceConfig{
	Nodes:              []string{"127.0.0.1:11087"},
	LogBucket:          "log-entries",
	IndexBucket:        "log-indexes",
	DaysToKeep:         30,
	EnableDebugLogging: true,
	RestPort:           8080,
}

func ReadConfig(path string) (ServiceConfig, error) {
	var config ServiceConfig
	if file, err := ioutil.ReadFile(path); err != nil {
		fmt.Printf("Error reading config file %q, %v\n", path, err)
		return config, err
	} else {
		if err := json.Unmarshal(file, &config); err != nil {
			fmt.Printf("Error parsing config file %v", err)
			return config, err
		}
	}
	return config, nil
}
