package riaktimeline

import (
	riak "github.com/basho/riak-go-client"
)

func buildNodeOptions(addresses []string) []*riak.NodeOptions {
	options := make([]*riak.NodeOptions, 0, len(addresses))
	for _, address := range addresses {
		options = append(options, &riak.NodeOptions{RemoteAddress: address})
	}
	return options
}

func initCluster(options []*riak.NodeOptions) (*riak.Cluster, err) {
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
