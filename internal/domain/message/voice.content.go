package message

import "github.com/mauriciorobertodev/whappy-go/internal/domain/file"

type VoiceContent struct {
	Voice    file.VoiceFile `json:"voice"`
	ViewOnce *bool          `json:"view_once"`
}

func NewVoiceContent(voice file.VoiceFile, viewOnce *bool) *VoiceContent {
	return &VoiceContent{
		Voice:    voice,
		ViewOnce: viewOnce,
	}
}

func (i *VoiceContent) Kind() MessageKind {
	return MessageKindVoice
}

func (i *VoiceContent) HasDuration() bool {
	return i.Voice.Duration != nil && *i.Voice.Duration > 0
}
