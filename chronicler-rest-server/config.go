package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServiceConfig struct {
	DaysToKeep         int           `json:"days-to-keep"`
	EnableDebugLogging bool          `json:"enable-debug"`
	RestPort           int           `json:"rest-port"`
	StorageConfig      StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Nodes        []string       `json:"nodes"`
	NodeTemplate RiakNodeConfig `json: node-template`
	LogBucket    string         `json:"log-bucket"`
	IndexBucket  string         `json:"index-bucket"`
}

type RiakNodeConfig struct {
	MinConnections uint16 `json:"min-connections"`
	MaxConnections uint16 `json:"max-connections"`
}

var defaultConfig = ServiceConfig{
	DaysToKeep:         30,
	EnableDebugLogging: true,
	RestPort:           8080,
	StorageConfig: StorageConfig{
		Nodes: []string{"127.0.0.1:11087"},
		NodeTemplate: RiakNodeConfig{
			MinConnections: 10,
			MaxConnections: 30,
		},
		LogBucket:   "log-entries",
		IndexBucket: "log-indexes",
	},
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
