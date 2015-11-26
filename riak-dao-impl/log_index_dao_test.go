package riakdaoimpl

import (
	"fmt"
	ent "github.com/PeerioTechnologies/chronicler/entity"
	"testing"
	"time"
)

func TestFetchAndAppendIndexEntry(t *testing.T) {
	initTestCluster()
	t.Log("Test cluster ", TestCluster, logIndexBucket)
	dao := NewLogIndexRiakDAO(TestCluster, logIndexBucket)
	millis := (time.Now().Nanosecond() % 1e6)
	testId := fmt.Sprintf("millis#%d", millis)
	entry := ent.IndexEntry{
		Key:     testId,
		Time:    time.Now(),
		Level:   "INFO",
		Type:    "Login",
		Caption: "This is first entry",
	}
	var before, after ent.TimelineIndex
	before, _ = dao.GetTimeline("testUser")
	if err := dao.AppendToTimeline("testUser", entry); err != nil {
		t.Errorf("Error saving entry", err)
	}
	after, _ = dao.GetTimeline("testUser")
	if len(before)+1 != len(after) {
		t.Errorf("expected length of timeline %v, got %v", len(before)+1, len(after))
	}
}
