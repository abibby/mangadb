package events

import (
	"gosalusa.com/event"
	"gosalusa.com/event/cron"
)

type FetchImagesEvent struct {
	cron.CronEvent

	Source string
	ID     string
}

var _ event.Event = (*FetchImagesEvent)(nil)

func (e *FetchImagesEvent) Type() event.EventType {
	return "mangadb:fetch-images"
}
