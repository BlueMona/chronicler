package entity

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	jsonString := "[{\"key\": \"aaa\", \"time\": \"2014-05-16T08:28:06.801064-04:00\", \"level\": \"ERROR\", \"type\": \"login\", \"caption\": \"Login by user anri\"},{\"key\": \"bbb\", \"time\": \"2015-07-28T13:34:18.801064-04:00\", \"level\": \"INFO\", \"type\": \"session\", \"caption\": \"Session by user anri\"},{\"key\": \"ccc\", \"time\": \"2015-01-08T13:04:01.801064-04:00\", \"level\": \"ERROR\", \"type\": \"login\", \"caption\": \"Login by user anri\"}]"
	var index TimelineIndex
	err := json.Unmarshal([]byte(jsonString), &index)
	if err != nil {
		t.Error(err)
	}
	if expected, actual := 3, len(index); expected != actual {
		t.Errorf("expected decoded Timeline index %v, got %v", expected, actual)
	}
	if expected, actual := "aaa", index[0].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
	if expected, actual := "ccc", index[2].Key; expected != actual {
		t.Errorf("expected key %v, got %v", expected, actual)
	}
}

func TestEncode(t *testing.T) {
	index := TimelineIndex{
		IndexEntry{
			Key:     "aaa",
			Time:    time.Now(),
			Level:   "INFO",
			Type:    "Login",
			Caption: "This is first entry",
		},
		IndexEntry{
			Key:     "bbb",
			Time:    time.Now(),
			Level:   "ERROR",
			Type:    "Login",
			Caption: "This is second entry",
		},
		IndexEntry{
			Key:     "ccc",
			Time:    time.Now(),
			Level:   "DEBUG",
			Type:    "Login",
			Caption: "This is third entry",
		},
	}
	output, err := json.Marshal(index)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(output))
}
