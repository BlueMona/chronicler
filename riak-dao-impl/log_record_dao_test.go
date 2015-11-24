package riakdaoimpl

import (
	"testing"
)

func TestSaveAndFetchLogRecord(t *testing.T) {
	initTestCluster()
	dao := NewLogRecordRiakDao(TestCluster, logRecordsBucket)

	msg := "log message"
	id := "1111111111"
	err1 := dao.SaveLogRecord(id, msg)
	if err1 != nil {
		t.Errorf("Error saving log", err1)
	}

	storedMsg, err2 := dao.GetLogRecord(id)
	if err2 != nil {
		t.Errorf("Error fetching log", err2)
	}

	if storedMsg != msg {
		t.Errorf("expected length of timeline \"%s\", got \"%s\"", msg, storedMsg)
	}

	dao.DeleteLogRecord(id)

	if storedMsg, _ := dao.GetLogRecord(id); storedMsg != "" {
		t.Errorf("Error deleting log. Value for key \"%s\" is \"%s\"", id, storedMsg)
	}
}
