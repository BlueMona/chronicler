package riakdaoimpl

import (
	riak "github.com/basho/riak-go-client"
)

func StartCluster(addresses []string, logIndexBucket string, enableDebug bool) (*riak.Cluster, error) {
	riak.EnableDebugLogging = enableDebug
	options := buildNodeOptions(addresses)
	cluster, err := initCluster(options)
	if err != nil {
		return cluster, err
	}
	if err := cluster.Start(); err != nil {
		return cluster, err
	}
	InitIndexBucket(cluster, logIndexBucket)
	return cluster, nil
}

func buildNodeOptions(addresses []string) []*riak.NodeOptions {
	options := make([]*riak.NodeOptions, 0, len(addresses))
	for _, address := range addresses {
		options = append(options, &riak.NodeOptions{RemoteAddress: address})
	}
	return options
}

func initCluster(options []*riak.NodeOptions) (*riak.Cluster, error) {
	nodes := make([]*riak.Node, 0, len(options))
	for _, nodeOpts := range options {
		if node, err := riak.NewNode(nodeOpts); err == nil {
			nodes = append(nodes, node)
		} else {
			logErr("Error initialising new Riak Node", err)
		}
	}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}
	return riak.NewCluster(opts)
}

func closeCluster(cluster *riak.Cluster) {
	if err := cluster.Stop(); err != nil {
		logErr("Error disconnecting from Riak cluster", err)
	}
}
