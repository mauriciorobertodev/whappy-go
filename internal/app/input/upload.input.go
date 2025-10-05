package input

import (
	"io"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
)

type UploadFile struct {
	Stream      io.Reader
	Metadata    file.Metadata
	ThumbnailID *string
}

type UpdateUploadMetadata struct {
	FileID   string
	Metadata file.Metadata
}

type ListUploads struct {
	Cursor *string
	Limit  int
}

func (i *ListUploads) Normalize() {
	if i.Limit <= 0 {
		i.Limit = 20
	}

	if i.Limit > 100 {
		i.Limit = 100
	}
}

type GetUpload struct {
	FileID string
}

type DeleteUpload struct {
	FileID string
}
