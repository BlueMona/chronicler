package riaktimeline

import (
	"sort"
	"time"
)

type IndexEntry struct {
	Key     string    `json:"key"`
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Type    string    `json:"type"`
	Caption string    `json:"caption"`
}

// implements sort.Interface for []IndexEntry based on Key field
type TimelineIndex []IndexEntry

func (tli TimelineIndex) Len() int           { return len(tli) }
func (tli TimelineIndex) Less(i, j int) bool { return tli[i].Key < tli[j].Key }
func (tli TimelineIndex) Swap(i, j int)      { tli[i], tli[j] = tli[j], tli[i] }

func removeDuplicateEntries(entries TimelineIndex) TimelineIndex {
	encountered := map[string]bool{}
	result := entries[:0]
	for _, entry := range entries {
		if encountered[entry.Key] != true {
			encountered[entry.Key] = true
			result = append(result, entry)
		}
	}
	return result
}

func splitByDaysAge(entries TimelineIndex, dayToLive int) (TimelineIndex, TimelineIndex) {
	newEntries := entries[:0]
	oldEntires := make(TimelineIndex, 0, len(entries))
	for _, entry := range entries {
		//TODO switch to calendar days difference
		if time.Since(entry.Time).Hours()/24 > float64(dayToLive) {
			oldEntires = append(oldEntires, entry)
		} else {
			newEntries = append(newEntries, entry)
		}
	}
	return newEntries, oldEntires
}

func sortEntries(entries TimelineIndex) TimelineIndex {
	sort.Sort(entries)
	return entries
}

func extractKeys(entries TimelineIndex) []string {
	result := make([]string, 0, len(entries))
	for _, entry := range entries {
		result = append(result, entry.Key)
	}
	return result
}

func mergeTimelines(list []TimelineIndex) TimelineIndex {
	merged := TimelineIndex{}
	for _, index := range list {
		merged = append(merged, index...)
	}
	return sortEntries(removeDuplicateEntries(merged))
}

func mergeAndCutOld(list []TimelineIndex) (TimelineIndex, []string) {
	merged := TimelineIndex{}
	for _, index := range list {
		merged = append(merged, index...)
	}
	merged = removeDuplicateEntries(merged)
	fresh, old := splitByDaysAge(merged, 30)
	oldIds := extractKeys(old)
	fresh = sortEntries(fresh)
	return fresh, oldIds
}
