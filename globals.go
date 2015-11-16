package riaktimeline

import (
	riak "github.com/basho/riak-go-client"
)

var Config TimelineConfig = defaultConfig
var RiakCluster *riak.Cluster
