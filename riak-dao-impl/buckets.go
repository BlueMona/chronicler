package riakdaoimpl

import (
	riak "github.com/basho/riak-go-client"
)

func initIndexBucket(cluster *riak.Cluster, bucket string) {
	cmd, err := riak.NewStoreBucketPropsCommandBuilder().
		WithBucket(bucket).
		WithAllowMult(true).
		Build()
	if err != nil {
		logError("initIndexBucket", "Error creating command for setting up bucket for \"allow_mult\"", err.Error())
	}
	if err = cluster.Execute(cmd); err != nil {
		logError("initIndexBucket", "Error setting up bucket for \"allow_mult\"", err.Error())
	}
}
