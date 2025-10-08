package input_test

import (
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/group"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group Inputs", func() {
	Describe("JoinGroup Input", func() {
		It("should validate successfully", func() {
			inp := &input.JoinGroup{
				Invite: "valid-invite",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when invite is empty", func() {
			inp := &input.JoinGroup{
				Invite: "",
			}
			Expect(inp.Validate()).ToNot(BeNil())
		})
	})

	Describe("LeaveGroup Input", func() {
		It("should validate successfully", func() {
			inp := &input.LeaveGroup{
				JID: "valid-jid",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.LeaveGroup{
				JID: "",
			}
			Expect(inp.Validate()).ToNot(BeNil())
		})
	})

	Describe("GetGroupInvite Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetGroupInvite{
				JID:     "valid-jid",
				Refresh: true,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.GetGroupInvite{
				JID:     "",
				Refresh: false,
			}
			Expect(inp.Validate()).ToNot(BeNil())
		})
	})

	Describe("UpdateGroupName Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupName{
				JID:  "valid-jid",
				Name: "New Group Name",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupName{
				JID:  "",
				Name: "New Group Name",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Name is empty", func() {
			inp := &input.UpdateGroupName{
				JID:  "valid-jid",
				Name: "",
			}
			Expect(inp.Validate()).To(Equal(group.ErrNameInvalid))
		})

		It("should fail validation when Name is too long", func() {
			longName := ""
			for i := 0; i < group.MaxNameLength+1; i++ {
				longName += "a"
			}
			inp := &input.UpdateGroupName{
				JID:  "valid-jid",
				Name: longName,
			}
			Expect(inp.Validate()).To(Equal(group.ErrNameTooLong))
		})
	})

	Describe("UpdateGroupDescription Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupDescription{
				JID:         "valid-jid",
				Description: "New Group Description",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupDescription{
				JID:         "",
				Description: "New Group Description",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Description is too long", func() {
			longDescription := ""
			for i := 0; i < group.MaxDescriptionLength+1; i++ {
				longDescription += "a"
			}
			inp := &input.UpdateGroupDescription{
				JID:         "valid-jid",
				Description: longDescription,
			}
			Expect(inp.Validate()).To(Equal(group.ErrDescriptionTooLong))
		})
	})

	Describe("UpdateGroupSetting Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupSetting{
				JID:     "valid-jid",
				Setting: group.GroupSettingAddParticipants,
				Policy:  group.GroupSettingPolicyAnyone,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupSetting{
				JID:     "",
				Setting: group.GroupSettingAddParticipants,
				Policy:  group.GroupSettingPolicyAnyone,
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Setting is invalid", func() {
			inp := &input.UpdateGroupSetting{
				JID:     "valid-jid",
				Setting: "invalid-setting",
				Policy:  group.GroupSettingPolicyAnyone,
			}
			Expect(inp.Validate()).To(Equal(group.ErrSettingInvalid))
		})

		It("should fail validation when Policy is invalid", func() {
			inp := &input.UpdateGroupSetting{
				JID:     "valid-jid",
				Setting: group.GroupSettingAddParticipants,
				Policy:  "invalid-policy",
			}
			Expect(inp.Validate()).To(Equal(group.ErrPolicyInvalid))
		})
	})

	Describe("UpdateGroupMessageDuration Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupMessageDuration{
				JID:      "valid-jid",
				Duration: group.GroupMessageDuration90Days,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupMessageDuration{
				JID:      "",
				Duration: group.GroupMessageDuration90Days,
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Duration is invalid", func() {
			inp := &input.UpdateGroupMessageDuration{
				JID:      "valid-jid",
				Duration: "invalid-duration",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidMessageDuration))
		})
	})

	Describe("UpdateGroupParticipants Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupParticipants{
				JID:          "valid-jid",
				Participants: []string{"participant-jid"},
				Action:       group.ParticipantsActionAdd,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupParticipants{
				JID:          "",
				Participants: []string{"participant-jid"},
				Action:       group.ParticipantsActionAdd,
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Participants is empty", func() {
			inp := &input.UpdateGroupParticipants{
				JID:          "valid-jid",
				Participants: []string{},
				Action:       group.ParticipantsActionAdd,
			}
			Expect(inp.Validate()).To(Equal(group.ErrRequireParticipants))
		})

		It("should fail validation when Action is invalid", func() {
			inp := &input.UpdateGroupParticipants{
				JID:          "valid-jid",
				Participants: []string{"participant-jid"},
				Action:       "invalid-action",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidAction))
		})
	})

	Describe("UpdateGroupPhoto Input", func() {
		It("should validate successfully", func() {
			inp := &input.UpdateGroupPhoto{
				JID:   "valid-jid",
				Photo: "base64-photo-string",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.UpdateGroupPhoto{
				JID:   "",
				Photo: "base64-photo-string",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})

		It("should fail validation when Photo is empty", func() {
			inp := &input.UpdateGroupPhoto{
				JID:   "valid-jid",
				Photo: "",
			}
			Expect(inp.Validate()).To(Equal(group.ErrPhotoRejected))
		})
	})

	Describe("RemoveGroupPhoto Input", func() {
		It("should validate successfully", func() {
			inp := &input.RemoveGroupPhoto{
				JID: "valid-jid",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.RemoveGroupPhoto{
				JID: "",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})
	})

	Describe("GetPhotoURL Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetPhotoURL{
				JID: "valid-jid",
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.GetPhotoURL{
				JID: "",
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})
	})

	Describe("CreateGroup Input", func() {
		It("should validate successfully", func() {
			inp := &input.CreateGroup{
				Name:         "New Group",
				Participants: []string{"participant1-jid", "participant2-jid"},
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when Name is empty", func() {
			inp := &input.CreateGroup{
				Name:         "",
				Participants: []string{"participant1-jid", "participant2-jid"},
			}
			Expect(inp.Validate()).To(Equal(group.ErrNameInvalid))
		})

		It("should fail validation when Name is too long", func() {
			longName := ""
			for i := 0; i < group.MaxNameLength+1; i++ {
				longName += "a"
			}
			inp := &input.CreateGroup{
				Name:         longName,
				Participants: []string{"participant1-jid", "participant2-jid"},
			}
			Expect(inp.Validate()).To(Equal(group.ErrNameTooLong))
		})
	})

	Describe("GetGroup Input", func() {
		It("should validate successfully", func() {
			inp := &input.GetGroup{
				JID:              "valid-jid",
				WithParticipants: nil,
			}
			Expect(inp.Validate()).To(BeNil())
		})

		It("should fail validation when JID is empty", func() {
			inp := &input.GetGroup{
				JID:              "",
				WithParticipants: nil,
			}
			Expect(inp.Validate()).To(Equal(group.ErrInvalidJID))
		})
	})
})
