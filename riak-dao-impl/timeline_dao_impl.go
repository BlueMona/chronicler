package riakdaoimpl

import (
	ent "github.com/PeerioTechnologies/riak-timeline-service/entity"
	riak "github.com/basho/riak-go-client"
	gosnow "github.com/sdming/gosnow"
	"strconv"
	"sync"
	"time"
)

//Implements TimelineDAO
type TimelineRiakDaoImpl struct {
	Cluster      *riak.Cluster
	IndexDao     *LogIndexRiakDAO
	LogRecordDao *LogRecordRiakDao
	SnowFlake    *gosnow.SnowFlake
}

func (dao *TimelineRiakDaoImpl) GetTimeline(id string, daysToKeep int) (ent.TimelineIndex, error) {
	index, err := dao.IndexDao.GetTimeline(id)
	if err != nil {
		return nil, err
	}
	//TODO add removal of old entries by channel
	index, _ = ent.SplitByDaysAge(index, daysToKeep)
	//async fill IndexEntry.Caption with real log data
	itemsQty := len(index)
	group := new(sync.WaitGroup)
	group.Add(itemsQty)
	for i := 0; i < itemsQty; i++ {
		go func(entry *ent.IndexEntry, group *sync.WaitGroup) {
			storedMsg, _ := dao.LogRecordDao.GetLogRecord(entry.Key)
			entry.Caption = storedMsg
			group.Done()
		}(&index[i], group)
	}
	group.Wait()
	return index, nil
}

func (dao *TimelineRiakDaoImpl) SaveLog(userId string, level string, typeStr string, msg string) error {
	id, _ := dao.SnowFlake.Next()
	idStr := strconv.FormatUint(id, 10)
	entry := ent.IndexEntry{
		Key:     idStr,
		Time:    time.Now(),
		Level:   level,
		Type:    typeStr,
		Caption: "",
	}
	errc := make(chan error)
	go func() {
		errc <- dao.LogRecordDao.SaveLogRecord(idStr, msg)
	}()
	go func() {
		errc <- dao.IndexDao.AppendToTimeline(userId, entry)
	}()
	var err error
	for i := 0; i < 2; i++ {
		if e := <-errc; e != nil {
			err = e
		}
	}
	return err
}

func NewTimelineRiakDaoImpl(cluster *riak.Cluster, indexBucket string, logBucket string) TimelineRiakDaoImpl {
	snowFlake, _ := gosnow.Default()
	return TimelineRiakDaoImpl{
		Cluster:      cluster,
		IndexDao:     NewLogIndexRiakDAO(cluster, indexBucket),
		LogRecordDao: NewLogRecordRiakDao(cluster, logBucket),
		SnowFlake:    snowFlake,
	}
}