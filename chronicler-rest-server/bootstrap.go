package main

import (
	"fmt"
	timeline "github.com/PeerioTechnologies/chronicler"
	riakDAO "github.com/PeerioTechnologies/chronicler/riak-dao-impl"
)

var config ServiceConfig
var dao timeline.TimelineDAO

func bootstrap(path string) {
	if cnf, err := ReadConfig(path); err != nil {
		config = defaultConfig
	} else {
		config = cnf
	}
	if cluster, err := riakDAO.StartCluster(
		config.StorageConfig.Nodes,
		config.StorageConfig.NodeTemplate.MinConnections,
		config.StorageConfig.NodeTemplate.MaxConnections,
		config.StorageConfig.IndexBucket,
		config.EnableDebugLogging); err == nil {
		dao = riakDAO.NewTimelineRiakDaoImpl(cluster, config.StorageConfig.IndexBucket, config.StorageConfig.LogBucket)
	} else {
		panic(fmt.Sprintf("Error initialising Riak DAO %s", err.Error()))
	}

}
