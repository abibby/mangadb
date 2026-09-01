package events

import (
	"abibby.com/salusa/event"
)

type MangaPlusSeries struct {
	ID string
}

var _ event.Event = (*MangaPlusSeries)(nil)

func (e *MangaPlusSeries) Type() event.EventType {
	return "icbmdb:mangaplus-series"
}
