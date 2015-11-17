package riaktimeline

import (
	"testing"
)

func TestSaveAndFetchLogRecord(t *testing.T) {
	Config = defaultConfig
	initTestCluster()
	msg := "log message"
	riakId, err1 := storeLog(msg)
	if err1 != nil {
		t.Errorf("Error saving log", err1)
	}

	storedMsg, err2 := fetchLog(riakId)
	if err2 != nil {
		t.Errorf("Error fetching log", err2)
	}

	if storedMsg != msg {
		t.Errorf("expected length of timeline \"%s\", got \"%s\"", msg, storedMsg)
	}

	deleteLog(riakId)

	if storedMsg, _ := fetchLog(riakId); storedMsg != "" {
		t.Errorf("Error deleting log. Value for key \"%s\" is \"%s\"", riakId, storedMsg)
	}
}
