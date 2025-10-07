package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

func NewEventInstanceCreated() events.Event {
	uuid := uuid.NewString()
	return events.New(
		instance.EventCreated,
		instance.PayloadInstanceCreated{
			ID:        uuid,
			Name:      "Fake Instance " + uuid,
			CreatedAt: time.Now().UTC(),
		},
		&uuid,
	)
}
