package riakdaoimpl

import (
	"encoding/json"
	ent "github.com/PeerioTechnologies/riak-timeline-service/entity"
	riak "github.com/basho/riak-go-client"
)

type LogIndexRiakDAO struct {
	Cluster     riak.Cluster
	Resolver    *TimelineConfilctResolver
	IndexBucket string
}

var resolver = &TimelineConfilctResolver{}

func (dao *LogIndexRiakDAO) fetch(userId string) *riak.FetchValueResponse {
	cmd, _ := riak.NewFetchValueCommandBuilder().
		WithBucket(dao.IndexBucket).
		WithNotFoundOk(true).
		WithKey(userId).
		WithConflictResolver(dao.Resolver).
		Build()
	if err := dao.Cluster.Execute(cmd); err != nil {
		logErr("Fetching timeline for "+userId, err)
		return nil
	}
	fvc := cmd.(*riak.FetchValueCommand)
	return fvc.Response
}

func (dao *LogIndexRiakDAO) AppendToTimeline(userId string, entry IndexEntry) error {
	value := riak.Object{}
	value.ContentType = "application/json"
	value.Charset = "utf-8"
	value.ContentEncoding = "utf-8"
	index := TimelineIndex{entry}
	if responce := dao.fetch(userId); responce != nil && !responce.IsNotFound {
		value.VClock = responce.VClock
		if err := json.Unmarshal(responce.Values[0].Value, &index); err != nil {
			return err
		}
	}
	index = ent.sortEntries(append(index, entry))
	encoded, _ := json.Marshal(index)
	value.Value = encoded
	if cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(dao.IndexBucket).
		WithKey(userId).
		WithContent(&value).
		Build(); err != nil {
		logErr("Saving timeline for "+userId, err)
		return err
	}
	if err = dao.Cluster.Execute(cmd); err != nil {
		logErr("Saving timeline for "+userId, err)
		return err
	}
	return nil
}

func (dao *LogIndexRiakDAO) getTimeline(userId string) (ent.TimelineIndex, error) {
	index := ent.TimelineIndex{}
	if responce := dao.fetch(userId); responce == nil || responce.IsNotFound {
		return index, nil
	}
	if err := json.Unmarshal(responce.Values[0].Value, &index); err != nil {
		return nil, err
	}
	index = ent.SortEntries(index)
	return index, nil
}

func NewLogIndexRiakDAO(cluster riak.Cluster, indexBucket string) LogIndexRiakDAO {
	return LogIndexRiakDAO{
		Cluster:     cluster,
		Resolver:    &TimelineConfilctResolver{},
		IndexBucket: indexBucket,
	}
}
