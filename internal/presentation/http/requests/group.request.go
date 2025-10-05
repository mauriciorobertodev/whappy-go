package requests

import (
	"fmt"

	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	"github.com/mauriciorobertodev/whappy-go/internal/presentation/http"
)

type GroupUpdateNameRequest struct {
	Name string `json:"name"`
}

func (r *GroupUpdateNameRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()
	if r.Name == "" {
		bag.Add("name", "Name is required")
	}

	if r.Name != "" && len(r.Name) > group.MaxNameLength {
		bag.Add("name", fmt.Sprintf("Name must be less than %d characters", group.MaxNameLength))
	}

	return bag
}

type GroupUpdateDescriptionRequest struct {
	Description *string `json:"description"` // nullable, set to null to remove description
}

func (r *GroupUpdateDescriptionRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if r.Description != nil && len(*r.Description) > group.MaxDescriptionLength {
		bag.Add("description", fmt.Sprintf("Description must be less than %d characters", group.MaxDescriptionLength))
	}

	return bag
}

type GroupUpdateTopicRequest struct {
	Topic *string `json:"topic"` // nullable, set to null to remove topic
}

func (r *GroupUpdateTopicRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if r.Topic != nil && len(*r.Topic) > group.MaxTopicLength {
		bag.Add("topic", fmt.Sprintf("Topic must be less than %d characters", group.MaxTopicLength))
	}

	return bag
}

type GroupUpdateSettingRequest struct {
	Setting group.GroupSettingName   `json:"setting"`
	Policy  group.GroupSettingPolicy `json:"policy"`
}

func (r *GroupUpdateSettingRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if !r.Setting.IsValid() {
		bag.Add("setting", "Invalid setting")
	}

	if !r.Policy.IsValid() {
		bag.Add("policy", "Invalid policy")
	}

	return bag
}

type GroupUpdateMessageDurationRequest struct {
	Duration group.GroupMessageDuration `json:"duration"`
}

func (r *GroupUpdateMessageDurationRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if !r.Duration.IsValid() {
		bag.Add("duration", "Invalid duration")
	}

	return bag
}

type GroupUpdatePhotoRequest struct {
	Photo string `json:"photo"` // base64 encoded
}

func (r *GroupUpdatePhotoRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if r.Photo == "" {
		bag.Add("photo", "Photo is required")
	}

	return bag
}

type UpdateParticipantsRequest struct {
	Participants []string                 `json:"participants"`
	Action       group.ParticipantsAction `json:"action"`
}

func (r *UpdateParticipantsRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if len(r.Participants) == 0 {
		bag.Add("participants", "Participants are required")
	}

	if r.Action != group.ParticipantsActionAdd && r.Action != group.ParticipantsActionRemove && r.Action != group.ParticipantsActionPromote && r.Action != group.ParticipantsActionDemote {
		bag.Add("action", "Action must be one of: add, remove, promote, demote")
	}

	return bag
}

type CreateGroupRequest struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

func (r *CreateGroupRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if r.Name == "" {
		bag.Add("name", "Name is required")
	}

	if len(r.Name) > group.MaxNameLength {
		bag.Add("name", fmt.Sprintf("Name must be less than %d characters", group.MaxNameLength))
	}

	return bag
}

func (r *CreateGroupRequest) ToInput() input.CreateGroup {
	return input.CreateGroup{
		Name:         r.Name,
		Participants: r.Participants,
	}
}

type UpdateGroupNameRequest struct {
	Name string `json:"name"`
}

func (r *UpdateGroupNameRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()
	if r.Name == "" {
		bag.Add("name", "Name is required")
	}

	if r.Name != "" && len(r.Name) > group.MaxNameLength {
		bag.Add("name", fmt.Sprintf("Name must be less than %d characters", group.MaxNameLength))
	}

	return bag
}

func (r *UpdateGroupNameRequest) ToInput(groupJID string) input.UpdateGroupName {
	return input.UpdateGroupName{
		JID:  groupJID,
		Name: r.Name,
	}
}

type UpdateGroupDescriptionRequest struct {
	Description string `json:"description"` // nullable, set to null to remove description
}

func (r *UpdateGroupDescriptionRequest) Validate() *http.ErrorBag {
	bag := http.NewErrorBag()

	if len(r.Description) > group.MaxDescriptionLength {
		bag.Add("description", fmt.Sprintf("Description must be less than %d characters", group.MaxDescriptionLength))
	}

	return bag
}

func (r *UpdateGroupDescriptionRequest) ToInput(groupJID string) input.UpdateGroupDescription {
	return input.UpdateGroupDescription{
		JID:         groupJID,
		Description: r.Description,
	}
}
