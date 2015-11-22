package riaktimeline

import (
	"fmt"
	"testing"
	"time"
)

func TestFetchAndAppendIndexEntry(t *testing.T) {
	Config = defaultConfig
	initTestCluster()
	millis := (time.Now().Nanosecond() % 1e6)
	testId := fmt.Sprintf("millis#%d", millis)
	entry := IndexEntry{
		Key:     testId,
		Time:    time.Now(),
		Level:   "INFO",
		Type:    "Login",
		Caption: "This is first entry",
	}
	var before, after TimelineIndex
	before, _ = getTimeline("testUser")
	if err := appendToTimeline("testUser", entry); err != nil {
		t.Errorf("Error saving entry", err)
	}
	after, _ = getTimeline("testUser")
	if len(before)+1 != len(after) {
		t.Errorf("expected length of timeline %v, got %v", len(before)+1, len(after))
	}
}
