package message

import "github.com/mauriciorobertodev/whappy-go/internal/domain/file"

type DocumentContent struct {
	Document  file.File `json:"document"`
	Thumbnail *string   `json:"thumbnail"`
	Caption   *string   `json:"caption"`
	Mentions  *[]string `json:"mentions"`
}

func NewDocumentContent(doc file.File, thumbnail *string, caption *string, mentions *[]string) *DocumentContent {
	return &DocumentContent{
		Document:  doc,
		Thumbnail: thumbnail,
		Caption:   caption,
		Mentions:  mentions,
	}
}

func (i *DocumentContent) Kind() MessageKind {
	return MessageKindDocument
}

func (i *DocumentContent) HasThumbnail() bool {
	return i.Thumbnail != nil
}
