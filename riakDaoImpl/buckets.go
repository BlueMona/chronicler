package riakdaoimpl

import (
	riak "github.com/basho/riak-go-client"
)

func initIndexBucket(bucket string) {
	if cmd, err := riak.NewStoreBucketPropsCommandBuilder().
		WithBucket(bucket).
		WithAllowMult(true).
		Build(); err != nil {
		logError("initIndexBucket", "Error creating command for setting up bucket for \"allow_mult\"", err.Error())
	}
	if err = RiakCluster.Execute(cmd); err != nil {
		logError("initIndexBucket", "Error setting up bucket for \"allow_mult\"", err.Error())
	}
}
