package message

import "github.com/mauriciorobertodev/whappy-go/internal/domain/file"

type AudioContent struct {
	Audio file.AudioFile `json:"audio"`
}

func NewAudioContent(audio file.AudioFile) *AudioContent {
	return &AudioContent{
		Audio: audio,
	}
}

func (i *AudioContent) Kind() MessageKind {
	return MessageKindAudio
}

func (i *AudioContent) HasDuration() bool {
	return i.Audio.Duration != nil && *i.Audio.Duration > 0
}
