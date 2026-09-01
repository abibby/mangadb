package events

import "abibby.com/salusa/event"

type VizSeries struct {
	ID string
}

var _ event.Event = (*VizSeries)(nil)

func (e *VizSeries) Type() event.EventType {
	return "icbmdb:viz-series"
}
