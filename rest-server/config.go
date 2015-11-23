package main

import (
	"encoding/json"
	riakDAO "github.com/PeerioTechnologies/riak-timeline-service/riakDaoImpl"
	"io/ioutil"
)

type ServiceConfig struct {
	*riakDAO.TimelineDAOConfig
	RestPort int `json:rest-port`
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
		return config, err
	}
	if err = json.Unmarshal(file, &config); err != nil {
		return config, err
	}
	return config, nil
}
