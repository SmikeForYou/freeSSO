package messages

import (
	"time"
)

type Event struct {
	Name string
	Time time.Time
}

type EventBuilder struct {
	event *Event
}

func (b *EventBuilder) SetName(name string) *EventBuilder {
	b.event.Name = name
	return b
}
func (b *EventBuilder) SetTime(t time.Time) *EventBuilder {
	b.event.Time = t
	return b
}

func NewEventBuilder() *EventBuilder { return &EventBuilder{event: &Event{Time: time.Now()}} }

type Readable struct {
	A string
	B int64
}

func (r Readable) Fields() []interface{} {
	return []interface{}{&r.A, &r.B}
}
