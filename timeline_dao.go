package chronicler

import (
	entity "github.com/PeerioTechnologies/chronicler/entity"
)

type TimelineDAO interface {
	StateReporter
	GetTimeline(id string) (entity.TimelineIndex, error)
	SaveLog(userId string, level string, typeStr string, msg string) error
}
