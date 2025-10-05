package message

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
)

type VideoContent struct {
	Video     file.VideoFile `json:"video"`
	Thumbnail *string        `json:"thumbnail"`
	Caption   *string        `json:"caption"`
	Mentions  *[]string      `json:"mentions"`
	ViewOnce  *bool          `json:"view_once"`
}

func NewVideoContent(video file.VideoFile, thumbnail *string, caption *string, mentions *[]string, viewOnce *bool) *VideoContent {
	return &VideoContent{
		Video:     video,
		Thumbnail: thumbnail,
		Caption:   caption,
		Mentions:  mentions,
		ViewOnce:  viewOnce,
	}
}

func (i *VideoContent) Kind() MessageKind {
	return MessageKindVideo
}

func (i *VideoContent) HasThumbnail() bool {
	return i.Thumbnail != nil
}

func (i *VideoContent) HasCaption() bool {
	return i.Caption != nil && *i.Caption != ""
}

func (i *VideoContent) HasMentions() bool {
	return i.Mentions != nil && len(*i.Mentions) > 0
}

func (i *VideoContent) TextForMentions() string {
	if i.Caption == nil {
		return ""
	}
	return *i.Caption
}
