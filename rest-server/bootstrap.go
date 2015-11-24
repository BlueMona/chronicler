package main

import (
	"fmt"
	timeline "github.com/PeerioTechnologies/riak-timeline-service"
	riakDAO "github.com/PeerioTechnologies/riak-timeline-service/riak-dao-impl"
)

var config ServiceConfig
var dao timeline.TimelineDAO

func bootstrap(path string) {
	if cnf, err := ReadConfig(path); err != nil {
		config = defaultConfig
	} else {
		config = cnf
	}
	if cluster, err := riakDAO.StartCluster(config.Nodes, config.IndexBucket, config.EnableDebugLogging); err == nil {
		dao = riakDAO.NewTimelineRiakDaoImpl(cluster, config.IndexBucket, config.LogBucket)
	} else {
		panic(fmt.Sprintf("Error initialising Riak DAO %s", err.Error()))
	}

}
