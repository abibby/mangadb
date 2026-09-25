package events

import (
	"abibby.com/mangadb/services/datasource"
	"gosalusa.com/event"
)

type FetchSeriesEvent struct {
	Series []*datasource.Series
}

var _ event.Event = (*FetchSeriesEvent)(nil)

func (e *FetchSeriesEvent) Type() event.EventType {
	return "mangadb:fetch-series"
}
