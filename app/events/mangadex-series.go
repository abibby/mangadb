package events

import (
	"abibby.com/salusa/event"
)

type MangadexSeries struct {
	ID string
}

var _ event.Event = (*MangadexSeries)(nil)

func (e *MangadexSeries) Type() event.EventType {
	return "mangadb:mangadex-series"
}
