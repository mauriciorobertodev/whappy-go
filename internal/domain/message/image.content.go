package message

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
)

type ImageContent struct {
	Image     *file.ImageFile `json:"image"`
	Thumbnail *string         `json:"thumbnail"`
	Caption   *string         `json:"caption"`
	Mentions  *[]string       `json:"mentions"`
	ViewOnce  *bool           `json:"view_once"`
}

func NewImageContent(image *file.ImageFile, thumbnail *string, caption *string, mentions *[]string, viewOnce *bool) ImageContent {
	return ImageContent{
		Image:     image,
		Thumbnail: thumbnail,
		Caption:   caption,
		Mentions:  mentions,
		ViewOnce:  viewOnce,
	}
}

func (i ImageContent) Kind() MessageKind {
	return MessageKindImage
}

func (i *ImageContent) HasThumbnail() bool {
	return i.Thumbnail != nil
}

func (i *ImageContent) HasCaption() bool {
	return i.Caption != nil && *i.Caption != ""
}

func (i *ImageContent) HasMentions() bool {
	return i.Mentions != nil && len(*i.Mentions) > 0
}

func (i *ImageContent) TextForMentions() string {
	if i.Caption == nil {
		return ""
	}
	return *i.Caption
}
