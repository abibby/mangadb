package events

import "abibby.com/salusa/event"

type AnilistSeries struct {
	ID string
}

var _ event.Event = (*AnilistSeries)(nil)

func (e *AnilistSeries) Type() event.EventType {
	return "mangadb:anilist-series"
}
