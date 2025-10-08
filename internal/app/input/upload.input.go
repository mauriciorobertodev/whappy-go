package input

import (
	"io"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
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

func (i *UpdateUploadMetadata) Validate() error {
	if !utils.IsUUID(i.FileID) {
		return file.ErrInvalidFileID
	}
	return nil
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

func (i *GetUpload) Validate() error {
	if !utils.IsUUID(i.FileID) {
		return file.ErrInvalidFileID
	}
	return nil
}

type DeleteUpload struct {
	FileID string
}

func (i *DeleteUpload) Validate() error {
	if !utils.IsUUID(i.FileID) {
		return file.ErrInvalidFileID
	}
	return nil
}
