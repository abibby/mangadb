package events

import (
	"abibby.com/salusa/event"
)

type FetchSeriesEvent struct {
	Source string
	ID     string
}

var _ event.Event = (*FetchSeriesEvent)(nil)

func (e *FetchSeriesEvent) Type() event.EventType {
	return "mangadb:fetch-series"
}
