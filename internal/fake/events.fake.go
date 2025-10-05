package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

func NewEventInstanceCreated() events.Event {
	uuid := uuid.NewString()
	return events.Event{
		Name: instance.EventCreated,
		Payload: instance.PayloadInstanceCreated{
			ID:        uuid,
			Name:      "Fake Instance " + uuid,
			CreatedAt: time.Now().UTC(),
		},
		OccurredAt: time.Now().UTC(),
		InstanceID: nil,
	}
}
