package riaktimeline

import (
	"encoding/json"
	riak "github.com/basho/riak-go-client"
)

//TODO extract storage interface for possible other storage backend

type TimelineConfilctResolver struct {
}

func (cr *TimelineConfilctResolver) Resolve(objs []*riak.Object) []*riak.Object {
	indexes := []TimelineIndex{}
	for _, obj := range objs {
		var index TimelineIndex
		err := json.Unmarshal(obj.Value, &index)
		if err == nil {
			indexes = append(indexes, index)
		}
	}
	mergedTimeline := mergeTimelines(indexes)
	jsonData, _ := json.Marshal(mergedTimeline)
	resultObject := &riak.Object{
		ContentType:     "application/json",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           jsonData,
	}
	return []*riak.Object{resultObject}
}

var resolver = &TimelineConfilctResolver{}

func queryTimeline(userId string) *riak.FetchValueResponse {
	cmd, _ := riak.NewFetchValueCommandBuilder().
		WithBucket(Config.IndexBucket).
		WithNotFoundOk(true).
		WithKey(userId).
		WithConflictResolver(resolver).
		Build()
	if err := RiakCluster.Execute(cmd); err != nil {
		logErr("Fetching timeline for "+userId, err)
		return nil
	}
	fvc := cmd.(*riak.FetchValueCommand)
	return fvc.Response
}

func appendToTimeline(userId string, entry IndexEntry) error {
	fetch := queryTimeline(userId)
	var index TimelineIndex
	var value riak.Object
	if fetch == nil || fetch.IsNotFound {
		value = riak.Object{}
		index = TimelineIndex{}
	} else {
		value := fetch.Values[0]
		if err := json.Unmarshal(value.Value, &index); err != nil {
			return err
		}
		index = sortEntries(removeDuplicateEntries(append(index, entry)))
		value.VClock = fetch.VClock
	}
	encoded, _ := json.Marshal(index)
	value.ContentType = "application/json"
	value.Charset = "utf-8"
	value.ContentEncoding = "utf-8"
	value.Value = encoded
	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucket(Config.IndexBucket).
		WithKey(userId).
		WithContent(&value).
		Build()
	if err != nil {
		logErr("Saving timeline for "+userId, err)
		return err
	}
	if err := RiakCluster.Execute(cmd); err != nil {
		logErr("Saving timeline for "+userId, err)
		return err
	}
	return nil
}

func getTimeline(userId string) (TimelineIndex, error) {
	fetch := queryTimeline(userId)
	index := TimelineIndex{}
	if fetch == nil || fetch.IsNotFound {
		return index, nil
	}
	if err := json.Unmarshal(fetch.Values[0].Value, &index); err != nil {
		return nil, err
	}
	index = sortEntries(removeDuplicateEntries(index))
	return index, nil
}
