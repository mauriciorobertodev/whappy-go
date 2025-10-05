package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

type eventFactory struct {
	prototype events.Event
}

func NewEvent() *eventFactory {
	return &eventFactory{
		prototype: events.Event{
			Name:       "",
			Payload:    nil,
			OccurredAt: time.Now().UTC(),
			InstanceID: nil,
		},
	}
}

func (f *eventFactory) WithName(name events.EventName) *eventFactory {
	f.prototype.Name = name
	return f
}

func (f *eventFactory) WithPayload(payload any) *eventFactory {
	f.prototype.Payload = payload
	return f
}

func (f *eventFactory) Create() events.Event {
	if f.prototype.Name == "" {
		f.prototype.Name = events.EventName("fake.event." + uuid.NewString())
	}

	if f.prototype.Payload == nil {
		f.prototype.Payload = "fake-payload " + uuid.NewString()
	}

	return f.prototype
}
