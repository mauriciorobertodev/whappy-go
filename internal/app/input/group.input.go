package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/group"

type JoinGroup struct {
	Invite string `json:"invite"`
}

func (inp *JoinGroup) Validate() error {
	if inp.Invite == "" {
		return group.ErrInviteInvalid
	}

	return nil
}

type LeaveGroup struct {
	JID string `json:"jid"`
}

func (inp *LeaveGroup) Validate() error {
	if inp.JID == "" {
		return group.ErrInviteInvalid
	}

	return nil
}

type GetGroupInvite struct {
	JID     string `json:"jid"`
	Refresh bool   `json:"refresh"`
}

func (inp *GetGroupInvite) Validate() error {
	if inp.JID == "" {
		return group.ErrInviteInvalid
	}

	return nil
}

type UpdateGroupName struct {
	JID  string `json:"jid"`
	Name string `json:"name"`
}

func (inp *UpdateGroupName) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if inp.Name == "" {
		return group.ErrNameInvalid
	}

	if len(inp.Name) > group.MaxNameLength {
		return group.ErrNameTooLong
	}

	return nil
}

type UpdateGroupDescription struct {
	JID         string `json:"jid"`
	Description string `json:"description"`
}

func (inp *UpdateGroupDescription) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if len(inp.Description) > group.MaxDescriptionLength {
		return group.ErrDescriptionTooLong
	}

	return nil
}

type UpdateGroupSetting struct {
	JID     string                   `json:"jid"`
	Setting group.GroupSettingName   `json:"setting"`
	Policy  group.GroupSettingPolicy `json:"policy"`
}

func (inp *UpdateGroupSetting) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if !inp.Setting.IsValid() {
		return group.ErrSettingInvalid
	}

	if !inp.Policy.IsValid() {
		return group.ErrPolicyInvalid
	}

	return nil
}

type UpdateGroupMessageDuration struct {
	JID      string                     `json:"jid"`
	Duration group.GroupMessageDuration `json:"duration"`
}

func (inp *UpdateGroupMessageDuration) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if !inp.Duration.IsValid() {
		return group.ErrInvalidMessageDuration
	}

	return nil
}

type UpdateGroupParticipants struct {
	JID          string                   `json:"jid"`
	Participants []string                 `json:"participants"`
	Action       group.ParticipantsAction `json:"action"`
}

func (inp *UpdateGroupParticipants) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if len(inp.Participants) == 0 {
		return group.ErrRequireParticipants
	}

	if !inp.Action.IsValid() {
		return group.ErrInvalidAction
	}

	return nil
}

type UpdateGroupPhoto struct {
	JID   string `json:"jid"`
	Photo string `json:"photo"`
}

func (inp *UpdateGroupPhoto) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	if inp.Photo == "" {
		return group.ErrPhotoRejected
	}

	return nil
}

type RemoveGroupPhoto struct {
	JID string `json:"jid"`
}

func (inp *RemoveGroupPhoto) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	return nil
}

type GetPhotoURL struct {
	JID     string `json:"jid"`
	Preview bool   `json:"preview"`
}

func (inp *GetPhotoURL) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	return nil
}

type CreateGroup struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

func (inp *CreateGroup) Validate() error {
	if inp.Name == "" {
		return group.ErrNameInvalid
	}

	if len(inp.Name) > group.MaxNameLength {
		return group.ErrNameTooLong
	}

	if len(inp.Participants) == 0 {
		return group.ErrRequireParticipants
	}

	return nil
}

type GetGroup struct {
	JID              string `json:"jid"`
	WithParticipants *bool  `json:"with_participants"`
}

func (inp *GetGroup) Validate() error {
	if inp.JID == "" {
		return group.ErrInvalidJID
	}

	return nil
}

type GetGroups struct {
	WithParticipants *bool `json:"with_participants"`
}
