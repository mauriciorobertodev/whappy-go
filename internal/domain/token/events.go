package token

import "github.com/mauriciorobertodev/whappy-go/internal/domain/events"

const (
	EventRenewed events.EventName = "token:renewed"
)

func (t *Token) EventRenewed(masked bool) events.Event {
	token := t.GetRaw()

	if masked {
		token = t.GetMasked()
	}

	return events.Event{
		Name: EventRenewed,
		Payload: PayloadTokenRenewed{
			ID:     t.ID,
			Token:  token,
			Masked: masked,
		},
		InstanceID: &t.InstanceID,
	}
}
