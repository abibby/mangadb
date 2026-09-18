package events

import (
	"gosalusa.com/event"
	"gosalusa.com/event/cron"
)

type UpdateViewsEvent struct {
	cron.CronEvent
}

var _ event.Event = (*UpdateViewsEvent)(nil)

func (e *UpdateViewsEvent) Type() event.EventType {
	return "mangadb:update-views"
}
