package riakdaoimpl

import (
	ent "github.com/PeerioTechnologies/chronicler/entity"
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
	DaysToKeep   int
}

func (dao *TimelineRiakDaoImpl) GetTimeline(id string) (ent.TimelineIndex, error) {
	index, err := dao.IndexDao.GetLofIndex(id)
	if err != nil {
		return nil, err
	}

	//add removal of old entries by channel
	index, oldEntries := ent.SplitByDaysAge(index, dao.DaysToKeep)
	go dao.removeOldEntries(oldEntries, make(chan bool))

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

func (dao *TimelineRiakDaoImpl) removeOldEntries(entries ent.TimelineIndex, ch chan bool) {
	itemsQty := len(entries)
	group := new(sync.WaitGroup)
	group.Add(itemsQty)
	for i := 0; i < itemsQty; i++ {
		go func(logId string, group *sync.WaitGroup) {
			dao.LogRecordDao.DeleteLogRecord(logId)
			group.Done()
		}(entries[i].Key, group)
	}
	group.Wait()
	ch <- true
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
		errc <- dao.IndexDao.AppendToLogIndex(userId, entry)
	}()
	var err error
	for i := 0; i < 2; i++ {
		if e := <-errc; e != nil {
			err = e
		}
	}
	return err
}

func NewTimelineRiakDaoImpl(cluster *riak.Cluster, indexBucket string, logBucket string, daysToKeep int) *TimelineRiakDaoImpl {
	snowFlake, _ := gosnow.Default()
	return &TimelineRiakDaoImpl{
		Cluster:      cluster,
		IndexDao:     NewLogIndexRiakDAO(cluster, indexBucket),
		LogRecordDao: NewLogRecordRiakDao(cluster, logBucket),
		SnowFlake:    snowFlake,
		DaysToKeep:   daysToKeep,
	}
}
