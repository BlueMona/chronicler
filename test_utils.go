package riaktimeline

import (
	"fmt"
	riak "github.com/basho/riak-go-client"
)

func initTestCluster() {
	riak.EnableDebugLogging = true
	EnableDebugLogging = true
	nodeOptions := buildNodeOptions(Config.Nodes)
	if RiakCluster == nil {
		RiakCluster = initCluster(nodeOptions)
		if err := RiakCluster.Start(); err != nil {
			panic(fmt.Sprintf("Error %v", err))
		}
		initIndexBucket(Config.IndexBucket)
	}
}
