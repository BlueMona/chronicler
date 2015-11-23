package riakdaoimpl

import (
	"encoding/json"
	ent "github.com/PeerioTechnologies/riak-timeline-service/entity"
	riak "github.com/basho/riak-go-client"
)

type TimelineConfilctResolver struct {
}

func (cr *TimelineConfilctResolver) Resolve(objs []*riak.Object) []*riak.Object {
	indexes := []TimelineIndex{}
	for _, obj := range objs {
		var index TimelineIndex
		if err := json.Unmarshal(obj.Value, &index); err == nil {
			indexes = append(indexes, index)
		}
	}
	mergedTimeline := ent.RemoveDuplicateEntries(ent.MergeTimelines(indexes))
	jsonData, _ := json.Marshal(mergedTimeline)
	resultObject := &riak.Object{
		ContentType:     "application/json",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           jsonData,
	}
	return []*riak.Object{resultObject}
}
