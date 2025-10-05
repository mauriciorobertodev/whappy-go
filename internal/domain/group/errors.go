package group

import "errors"

var (
	ErrInvalidJID             = errors.New("invalid group JID")
	ErrInviteInvalid          = errors.New("invalid invite")
	ErrInviteRevoked          = errors.New("revoked invite")
	ErrInviteUnauthorized     = errors.New("not authorized to fetch group invite link")
	ErrNameInvalid            = errors.New("group name invalid")
	ErrNameTooLong            = errors.New("group name too long")
	ErrDescriptionInvalid     = errors.New("group description invalid")
	ErrDescriptionTooLong     = errors.New("group description too long")
	ErrTopicTooLong           = errors.New("group topic too long")
	ErrNotFound               = errors.New("group not found")
	ErrMaybeNotMember         = errors.New("maybe not a member of the group")
	ErrNotMember              = errors.New("user is not a member of the group")
	ErrSettingNotSupported    = errors.New("group setting not supported")
	ErrSettingInvalid         = errors.New("group setting invalid")
	ErrPolicyInvalid          = errors.New("group policy invalid")
	ErrInvalidGroupJID        = errors.New("invalid group JID")
	ErrInvalidMessageDuration = errors.New("invalid message duration")
	ErrInvalidAction          = errors.New("invalid action")
	ErrRequireParticipants    = errors.New("requires participants")

	ErrPhotoUnsupportedFormat = errors.New("only jpeg and png are allowed")
	ErrPhotoInvalidDimensions = errors.New("dimensions must be between 192x192 and 640x640")
	ErrPhotoRejected          = errors.New("photo rejected")
)
