package riaktimeline

import (
	"testing"
)

func TestSaveAndFetchLogRecord(t *testing.T) {
	Config = defaultConfig
	initTestCluster()
	msg := "log message"
	id := "1111111111"
	err1 := storeLog(id, msg)
	if err1 != nil {
		t.Errorf("Error saving log", err1)
	}

	storedMsg, err2 := fetchLog(id)
	if err2 != nil {
		t.Errorf("Error fetching log", err2)
	}

	if storedMsg != msg {
		t.Errorf("expected length of timeline \"%s\", got \"%s\"", msg, storedMsg)
	}

	deleteLog(id)

	if storedMsg, _ := fetchLog(id); storedMsg != "" {
		t.Errorf("Error deleting log. Value for key \"%s\" is \"%s\"", id, storedMsg)
	}
}
