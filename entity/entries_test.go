package entity

import (
	"testing"
	"time"
)

func TestSorting(t *testing.T) {
	index := TimelineIndex{
		IndexEntry{
			Key:     "aaa",
			Time:    time.Now(),
			Level:   "INFO",
			Type:    "Login",
			Caption: "This is first entry",
		},
		IndexEntry{
			Key:     "ccc",
			Time:    time.Now(),
			Level:   "ERROR",
			Type:    "Login",
			Caption: "This is second entry",
		},
		IndexEntry{
			Key:     "bbb",
			Time:    time.Now(),
			Level:   "DEBUG",
			Type:    "Login",
			Caption: "This is third entry",
		},
	}
	index = SortEntries(index)
	if expected, actual := "aaa", index[0].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
	if expected, actual := "bbb", index[1].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
	if expected, actual := "ccc", index[2].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
}

func TestRemovingDuplicates(t *testing.T) {
	index := TimelineIndex{
		IndexEntry{
			Key:     "aaa",
			Time:    time.Now(),
			Level:   "INFO",
			Type:    "Login",
			Caption: "This is first entry",
		},
		IndexEntry{
			Key:     "ccc",
			Time:    time.Now(),
			Level:   "ERROR",
			Type:    "Login",
			Caption: "This is second entry",
		},
		IndexEntry{
			Key:     "bbb",
			Time:    time.Now(),
			Level:   "DEBUG",
			Type:    "Login",
			Caption: "This is third entry",
		},
		IndexEntry{
			Key:     "aaa",
			Time:    time.Now(),
			Level:   "INFO",
			Type:    "Login",
			Caption: "This is first entry",
		},
		IndexEntry{
			Key:     "ccc",
			Time:    time.Now(),
			Level:   "ERROR",
			Type:    "Login",
			Caption: "This is second entry",
		},
		IndexEntry{
			Key:     "bbb",
			Time:    time.Now(),
			Level:   "DEBUG",
			Type:    "Login",
			Caption: "This is third entry",
		},
		IndexEntry{
			Key:     "aaa",
			Time:    time.Now(),
			Level:   "INFO",
			Type:    "Login",
			Caption: "This is first entry",
		},
		IndexEntry{
			Key:     "ccc",
			Time:    time.Now(),
			Level:   "ERROR",
			Type:    "Login",
			Caption: "This is second entry",
		},
		IndexEntry{
			Key:     "bbb",
			Time:    time.Now(),
			Level:   "DEBUG",
			Type:    "Login",
			Caption: "This is third entry",
		},
	}
	index = SortEntries(RemoveDuplicateEntries(index))
	if expected, actual := 3, len(index); expected != actual {
		t.Errorf("expected length %v, got %v", expected, actual)
	}
	if expected, actual := "aaa", index[0].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
	if expected, actual := "bbb", index[1].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
	if expected, actual := "ccc", index[2].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
}
