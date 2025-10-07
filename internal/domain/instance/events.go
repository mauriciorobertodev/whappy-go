package instance

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

// To listen all instance events use "instance:*"
const (
	// To listen all instance events, use the prefix "instance/*"
	EventCreated events.EventName = "instance:created"
	EventDeleted events.EventName = "instance:deleted"
	EventUpdated events.EventName = "instance:updated"
	EventToken   events.EventName = "instance:token"
	// To listen all session events, use the prefix "instance:session/*"
	EventSessionLoggedIn     events.EventName = "instance:session/logged_in"
	EventSessionLoggedOut    events.EventName = "instance:session/logged_out"
	EventSessionConnecting   events.EventName = "instance:session/connecting"
	EventSessionConnected    events.EventName = "instance:session/connected"
	EventSessionDisconnected events.EventName = "instance:session/disconnected"
	EventSessionError        events.EventName = "instance:session/error"
	// To listen all pairing events, use the prefix "instance:pairing/*"
	EventPairingStarted events.EventName = "instance:pairing/started"
	EventPairingQRCode  events.EventName = "instance:pairing/qr"
	EventPairingSuccess events.EventName = "instance:pairing/success"
	EventPairingFailed  events.EventName = "instance:pairing/failed"
)

// #region Pairing Events
func (i *Instance) EventPairingStarted() events.Event {
	return events.New(
		EventPairingStarted,
		PayloadInstancePairingStarted{
			ID:   i.ID,
			Name: i.Name,
		},
		&i.ID,
	)
}

func (i *Instance) EventQRCodeGenerated(qr string) events.Event {
	return events.New(
		EventPairingQRCode,
		PayloadInstanceQRCodeGenerated{
			ID:     i.ID,
			QRCode: qr,
		},
		&i.ID,
	)
}

func (i *Instance) EventPairingFailed(code FailPairingFailedCode, attemptedPhone string, err error) events.Event {
	return events.New(
		EventPairingFailed,
		PayloadInstancePairingFailed{
			ID:             i.ID,
			Phone:          i.Phone,
			Code:           code,
			AttemptedPhone: attemptedPhone,
			Error:          err,
		},
		&i.ID,
	)
}

// #region Login Events
func (i *Instance) EventLoggedIn() events.Event {
	return events.New(
		EventSessionLoggedIn,
		PayloadInstanceLoggedIn{
			ID:     i.ID,
			Name:   i.Name,
			Phone:  i.Phone,
			JID:    i.JID,
			Device: i.Device,
		},
		&i.ID,
	)
}

// #region Creation Events
func (i *Instance) EventCreated() events.Event {
	return events.New(
		EventCreated,
		PayloadInstanceCreated{
			ID:   i.ID,
			Name: i.Name,
		},
		&i.ID,
	)
}

// #region Connection Events
func (i *Instance) EventConnecting() events.Event {
	return events.New(
		EventSessionConnecting,
		PayloadInstanceConnecting{
			ID:   i.ID,
			Name: i.Name,
		},
		&i.ID,
	)
}

func (i *Instance) EventConnectionFailed(err string) events.Event {
	return events.New(
		EventSessionError,
		PayloadInstanceConnectionFailed{
			ID:    i.ID,
			Name:  i.Name,
			Error: err,
		},
		&i.ID,
	)
}

func (i *Instance) EventConnected() events.Event {
	return events.New(
		EventSessionConnected,
		PayloadInstanceConnected{
			ID:    i.ID,
			Name:  i.Name,
			Phone: i.Phone,
		},
		&i.ID,
	)
}

func (i *Instance) EventDisconnected() events.Event {
	return events.New(
		EventSessionDisconnected,
		PayloadInstanceConnected{
			ID:    i.ID,
			Name:  i.Name,
			Phone: i.Phone,
		},
		&i.ID,
	)
}

func (i *Instance) EventLoggedOut() events.Event {
	return events.New(
		EventSessionLoggedOut,
		PayloadInstanceLoggedOut{
			ID:   i.ID,
			Name: i.Name,
		},
		&i.ID,
	)
}
