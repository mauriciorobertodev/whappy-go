package group

import (
	"time"
)

type GroupSettingName string

const (
	GroupSettingSendMessages        GroupSettingName = "send_messages"
	GroupSettingEditGroupInfo       GroupSettingName = "edit_group_info"
	GroupSettingApproveParticipants GroupSettingName = "approve_participants"
	GroupSettingAddParticipants     GroupSettingName = "add_participants"
)

func (gsn GroupSettingName) IsValid() bool {
	switch gsn {
	case
		GroupSettingSendMessages,
		GroupSettingEditGroupInfo,
		GroupSettingApproveParticipants,
		GroupSettingAddParticipants:
		return true
	default:
		return false
	}
}

type GroupSettingPolicy string

const (
	GroupSettingPolicyAnyone GroupSettingPolicy = "anyone" // All members
	GroupSettingPolicyAdmins GroupSettingPolicy = "admins" // Just admins
)

func (gsp GroupSettingPolicy) IsValid() bool {
	switch gsp {
	case
		GroupSettingPolicyAnyone,
		GroupSettingPolicyAdmins:
		return true
	default:
		return false
	}
}

type GroupMessageDuration string

const (
	GroupMessageDurationOff     GroupMessageDuration = "off"
	GroupMessageDuration24Hours GroupMessageDuration = "24h"
	GroupMessageDuration7Days   GroupMessageDuration = "7d"
	GroupMessageDuration90Days  GroupMessageDuration = "90d"
	GroupMessageDurationCustom  GroupMessageDuration = "custom"
)

func (e GroupMessageDuration) IsValid() bool {
	switch e {
	case
		GroupMessageDurationOff,
		GroupMessageDuration24Hours,
		GroupMessageDuration7Days,
		GroupMessageDuration90Days,
		GroupMessageDurationCustom:
		return true
	default:
		return false
	}
}

func (d GroupMessageDuration) ToExpiration() uint32 {
	switch d {
	case GroupMessageDurationOff:
		return 0
	case GroupMessageDuration24Hours:
		return 24 * 60 * 60
	case GroupMessageDuration7Days:
		return 7 * 24 * 60 * 60
	case GroupMessageDuration90Days:
		return 90 * 24 * 60 * 60
	case GroupMessageDurationCustom:
		return 0 // Custom duration, should be set manually
	default:
		return 0
	}
}

func (d GroupMessageDuration) ToDuration() time.Duration {
	return time.Duration(d.ToExpiration()) * time.Second
}

func NewGroupMessageDurationFromSeconds(seconds uint32) GroupMessageDuration {
	switch seconds {
	case 0:
		return GroupMessageDurationOff
	case 24 * 60 * 60:
		return GroupMessageDuration24Hours
	case 7 * 24 * 60 * 60:
		return GroupMessageDuration7Days
	case 90 * 24 * 60 * 60:
		return GroupMessageDuration90Days
	default:
		return GroupMessageDurationCustom
	}
}

type ParticipantsAction string

const (
	ParticipantsActionAdd     ParticipantsAction = "add"
	ParticipantsActionRemove  ParticipantsAction = "remove"
	ParticipantsActionPromote ParticipantsAction = "promote"
	ParticipantsActionDemote  ParticipantsAction = "demote"
)

func (a ParticipantsAction) IsValid() bool {
	switch a {
	case
		ParticipantsActionAdd,
		ParticipantsActionRemove,
		ParticipantsActionPromote,
		ParticipantsActionDemote:
		return true
	default:
		return false
	}
}
