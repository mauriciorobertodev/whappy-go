package instance

type InstanceStatus string

const (
	StatusCreated           InstanceStatus = "CREATED"            // for new instances
	StatusPairing           InstanceStatus = "PAIRING"            // when waiting for QR code scan or pair code
	StatusLoggedIn          InstanceStatus = "LOGGED_IN"          // when successfully logged in but not connected to websocket
	StatusLoggedOut         InstanceStatus = "LOGGED_OUT"         // when logged out (needs to pair again)
	StatusConnecting        InstanceStatus = "CONNECTING"         // when trying to connect on websocket
	StatusConnected         InstanceStatus = "CONNECTED"          // when trying to reconnect on websocket
	StatusTemporaryBanned   InstanceStatus = "TEMPORARY_BANNED"   // when temporarily banned (e.g., too many requests, too many block by people, etc)
	StatusPermanentlyBanned InstanceStatus = "PERMANENTLY_BANNED" // when permanently banned (e.g., by WhatsApp)
)

// returns true if the status represents a new instance (never logged in)
func (s InstanceStatus) IsNew() bool {
	return s == StatusCreated
}

// returns true if the status represents a pairing state (waiting for QR code scan or pair code)
func (s InstanceStatus) IsPairing() bool {
	return s == StatusPairing
}

// returns true if the status represents a logged-in state (but not necessarily connected to websocket)
func (s InstanceStatus) IsLoggedIn() bool {
	return s == StatusLoggedIn || s == StatusConnected || s == StatusConnecting
}

// returns true if the status represents a logged-out state (needs to pair again) / its has been logged in before but now is logged out
func (s InstanceStatus) IsLoggedOut() bool {
	return s == StatusLoggedOut
}

// returns true if the status represents a connecting state to websocket, logged in but not yet connected to websocket
func (s InstanceStatus) IsConnecting() bool {
	return s == StatusConnecting
}

// returns true if the status represents a connected state to websocket
func (s InstanceStatus) IsConnected() bool {
	return s == StatusConnected
}

// returns true if the status represents a disconnected state from websocket
func (s InstanceStatus) IsDisconnected() bool {
	return s == StatusLoggedOut ||
		s == StatusCreated ||
		s == StatusTemporaryBanned ||
		s == StatusPermanentlyBanned
}

// returns true if the status represents an active state (logged in and connected to websocket)
func (s InstanceStatus) IsActive() bool {
	return s == StatusLoggedIn || s == StatusConnected
}

// returns true if the status represents a banned state (temporary or permanent)
func (s InstanceStatus) IsBanned() bool {
	return s == StatusTemporaryBanned || s == StatusPermanentlyBanned
}
