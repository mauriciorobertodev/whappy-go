package instance

import "errors"

var (
	ErrInstanceIsConnecting   = errors.New("instance is connecting")
	ErrInstanceIsPairing      = errors.New("instance is pairing")
	ErrInstanceIsBanned       = errors.New("instance is banned")
	ErrInstanceIsDisconnected = errors.New("instance is disconnected")
	ErrInstanceIsConnected    = errors.New("instance is connected")
	ErrInstanceIsLoggedOut    = errors.New("instance is logged out")

	ErrInstanceNotLoggedIn  = errors.New("instance is not logged in")
	ErrInstanceNotConnected = errors.New("instance is not connected")
	ErrInstanceNotPairing   = errors.New("instance is not pairing")

	ErrInstanceAlreadyPairing      = errors.New("instance is already pairing")
	ErrInstanceAlreadyLoggedIn     = errors.New("instance is already logged in")
	ErrInstanceAlreadyLoggedOut    = errors.New("instance is already logged out")
	ErrInstanceAlreadyConnecting   = errors.New("instance is already connecting")
	ErrInstanceAlreadyConnected    = errors.New("instance is already connected")
	ErrInstanceAlreadyDisconnected = errors.New("instance is already disconnected")
	ErrInstanceAlreadyPaired       = errors.New("instance is already paired")

	ErrInstanceBanned = errors.New("instance is banned")

	ErrInstancePhoneConflict = errors.New("instance phone conflict, the phone is different from the last logged in")
	ErrInstanceNotFound      = errors.New("instance not found")

	ErrInstanceNoQrCode = errors.New("instance has no QR code")

	ErrNameTooLong = errors.New("instance name is too long")
	ErrInvalidID   = errors.New("invalid instance ID")
)
