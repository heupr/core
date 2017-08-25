package replay

import (
	"core/pipeline/ingestor"
	"sync"
)

var eventPool = sync.Pool{
	New: func() interface{} {
		return &ingestor.Event{}
	},
}

// GetEvent gets Event struct from sync pool
func GetEvent() *ingestor.Event {
	return eventPool.Get().(*ingestor.Event)
}

// RecycleEvent puts back Event struct back into sync pool
func RecycleEvent(event *ingestor.Event) {
	eventPool.Put(event)
}

type EventList struct {
	Data  []ingestor.Event
	Count int
}

var eventListPool = sync.Pool{
	New: func() interface{} {
		return &EventList{}
	},
}

// GetEventList gets EventList struct from sync pool
func GetEventList() *EventList {
	return eventListPool.Get().(*EventList)
}

// RecycleEventList puts back EventList struct back into sync pool
func RecycleEventList(eventList *EventList) {
	eventListPool.Put(eventList)
}
