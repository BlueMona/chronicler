package riakdaoimpl

import (
	"fmt"
	riak "github.com/basho/riak-go-client"
)

var testRiakNodes []string = []string{"127.0.0.1:11087"}
var logIndexBucket string = "log-indexes"
var logRecordsBucket string = "log-records"
var TestCluster *riak.Cluster

func initTestCluster() {
	if TestCluster == nil {
		riak.EnableDebugLogging = true
		nodeOptions := buildNodeOptions(testRiakNodes)
		cluster, err := initCluster(nodeOptions)
		if err != nil {
			panic(fmt.Sprintf("Error %v", err))
		}
		if err := cluster.Start(); err != nil {
			panic(fmt.Sprintf("Error %v", err))
		}
		initIndexBucket(cluster, logIndexBucket)
		TestCluster = cluster
	}
}
