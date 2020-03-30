package util

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Record struct {
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  time.Duration     `json:"duration"`
	Tags      map[string]string `json:"tags"`
}

func (r Record) String() string {
	tags, err := json.Marshal(r.Tags)
	if err != nil {
		tags = []byte("screwed up")
	}
	return fmt.Sprintf("%s start=%s end=%s Duration=%s durnano=%d", tags, r.StartTime.UTC().String(), r.EndTime.UTC().String(), r.Duration.String(), r.Duration.Nanoseconds())

}

func mergeMaps(left, right map[string]string) map[string]string {
	newMap := map[string]string{}
	for k, v := range left {
		newMap[k] = v
	}
	for k, v := range right {
		newMap[k] = v
	}
	return newMap
}

type TimeKeeper struct {
	parent  *TimeKeeper
	tags    map[string]string
	records []Record
}

func (tk *TimeKeeper) Child(tags map[string]string) *TimeKeeper {
	return &TimeKeeper{
		parent: tk,
		tags:   tags,
	}
}

func (tk *TimeKeeper) Start() *Timer {
	timer := &Timer{
		tracker:   tk,
		startTime: time.Now(),
	}
	return timer
}

func (tk *TimeKeeper) Records() []Record {
	if tk.parent != nil {
		return tk.parent.Records()
	}
	return tk.records
}

func (tk *TimeKeeper) AddRecord(record Record) {
	record.Tags = mergeMaps(record.Tags, tk.tags)
	if tk.parent != nil {
		tk.parent.AddRecord(record)
	} else {
		tk.records = append(tk.records, record)
	}

}

type Timer struct {
	tracker   *TimeKeeper
	startTime time.Time
	stopTime  time.Time
	elapsed   time.Duration
	stopOnce  sync.Once
}

func (t *Timer) Stop() {
	stopTime := time.Now()
	t.stopOnce.Do(func() {
		t.stopTime = stopTime
		t.elapsed = t.stopTime.Sub(t.startTime)
	})
}

func (t *Timer) Save(localTags map[string]string) {
	t.Stop()
	t.tracker.AddRecord(Record{StartTime: t.startTime, EndTime: t.stopTime, Duration: t.elapsed, Tags: localTags})
}

func NewTimeKeeper(tags map[string]string) *TimeKeeper {
	return &TimeKeeper{
		tags:    tags,
		records: []Record{},
	}

}
