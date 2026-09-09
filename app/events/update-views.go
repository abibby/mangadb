package events

import (
	"abibby.com/salusa/event"
	"abibby.com/salusa/event/cron"
)

type UpdateViewsEvent struct {
	cron.CronEvent
}

var _ event.Event = (*UpdateViewsEvent)(nil)

func (e *UpdateViewsEvent) Type() event.EventType {
	return "mangadb:update-views"
}
