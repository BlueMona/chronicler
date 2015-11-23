package timeline

import (
	entity "github.com/PeerioTechnologies/riak-timeline-service/entity"
)

type TimelineDAO interface {
	GetTimeline(id string) (entity.TimelineIndex, error)
	SaveLog(userId string, level string, typeStr string, msg string) error
}
