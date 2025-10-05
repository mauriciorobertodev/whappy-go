package models

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
)

type SQLThumbnailFile struct {
	ID        *string `db:"id"`
	Name      *string `db:"name"`
	Mime      *string `db:"mime"`
	Extension *string `db:"extension"`

	Size   *uint64 `db:"size"`
	Sha256 *string `db:"sha256"`

	Path *string `db:"path"`
	URL  *string `db:"url"`

	Width    *uint32 `db:"width"`
	Height   *uint32 `db:"height"`
	Duration *uint32 `db:"duration"`
	Pages    *uint32 `db:"pages"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`

	InstanceID  *string `db:"instance_id"`
	ThumbnailID *string `db:"thumbnail_id"`
}

func (s *SQLThumbnailFile) ToEntity() (*file.ImageFile, error) {
	if s == nil {
		return nil, nil
	}

	return &file.ImageFile{
		File: file.File{
			ID:        *s.ID,
			Name:      *s.Name,
			Mime:      *s.Mime,
			Extension: *s.Extension,

			Size:   *s.Size,
			Sha256: *s.Sha256,

			Path: *s.Path,
			URL:  *s.URL,

			Width:    s.Width,
			Height:   s.Height,
			Duration: s.Duration,
			Pages:    s.Pages,

			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,

			InstanceID: s.InstanceID,
		},
		Width:  s.Width,
		Height: s.Height,
	}, nil
}

type SQLFile struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Mime      string `db:"mime"`
	Extension string `db:"extension"`

	Size   uint64 `db:"size"`
	Sha256 string `db:"sha256"`

	Path string `db:"path"`
	URL  string `db:"url"`

	Width    *uint32 `db:"width"`
	Height   *uint32 `db:"height"`
	Duration *uint32 `db:"duration"`
	Pages    *uint32 `db:"pages"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	InstanceID  *string `db:"instance_id"`
	ThumbnailID *string `db:"thumbnail_id"`

	Thumbnail *SQLThumbnailFile `db:"thumbnail"`
}

func (s *SQLFile) ToEntity() (*file.File, error) {
	var createdAt, updatedAt *time.Time
	tc := s.CreatedAt.UTC()
	createdAt = &tc
	tu := s.UpdatedAt.UTC()
	updatedAt = &tu

	var thumbnail *file.ImageFile
	if s.Thumbnail != nil && s.Thumbnail.ID != nil {
		var err error
		thumbnail, err = s.Thumbnail.ToEntity()
		if err != nil {
			return nil, err
		}
	}

	return &file.File{
		ID:        s.ID,
		Name:      s.Name,
		Mime:      s.Mime,
		Extension: s.Extension,

		Size:   s.Size,
		Sha256: s.Sha256,

		Path: s.Path,
		URL:  s.URL,

		Width:    s.Width,
		Height:   s.Height,
		Duration: s.Duration,
		Pages:    s.Pages,

		CreatedAt: createdAt,
		UpdatedAt: updatedAt,

		InstanceID: s.InstanceID,
		Thumbnail:  thumbnail,
	}, nil
}

func FromFileEntity(file *file.File) (*SQLFile, error) {
	var thumbnailID *string
	if file.Thumbnail != nil {
		thumbnailID = &file.Thumbnail.ID
	}

	var createdAt, updatedAt *time.Time
	if file.CreatedAt != nil {
		t := file.CreatedAt.UTC()
		createdAt = &t
	} else {
		now := time.Now().UTC()
		createdAt = &now
	}

	if file.UpdatedAt != nil {
		t := file.UpdatedAt.UTC()
		updatedAt = &t
	} else {
		now := time.Now().UTC()
		updatedAt = &now
	}

	return &SQLFile{
		ID:        file.ID,
		Name:      file.Name,
		Mime:      file.Mime,
		Extension: file.Extension,

		Size:   file.Size,
		Sha256: file.Sha256,

		Path: file.Path,
		URL:  file.URL,

		Width:    file.Width,
		Height:   file.Height,
		Duration: file.Duration,
		Pages:    file.Pages,

		CreatedAt: *createdAt,
		UpdatedAt: *updatedAt,

		InstanceID:  file.InstanceID,
		ThumbnailID: thumbnailID,
	}, nil
}
