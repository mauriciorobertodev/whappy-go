package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
)

type fileFactory struct {
	prototype *file.File
}

func FileFactory() *fileFactory {
	return &fileFactory{
		prototype: &file.File{
			ID:        "",
			Name:      "",
			Mime:      "",
			Extension: "",

			Size:      0,
			Sha256:    "",
			Sha256Enc: "",
			MediaKey:  "",

			Path:       "",
			DirectPath: "",
			URL:        "",

			Width:    nil,
			Height:   nil,
			Duration: nil,
			Pages:    nil,

			Thumbnail: nil,

			CreatedAt: nil,
			UpdatedAt: nil,

			InstanceID: nil,
		},
	}
}

func (f *fileFactory) WithID(id string) *fileFactory {
	f.prototype.ID = id
	return f
}

func (f *fileFactory) WithName(name string) *fileFactory {
	f.prototype.Name = name
	return f
}

func (f *fileFactory) WithMime(mime string) *fileFactory {
	f.prototype.Mime = mime
	return f
}

func (f *fileFactory) WithExtension(extension string) *fileFactory {
	f.prototype.Extension = extension
	return f
}

func (f *fileFactory) WithSize(size uint64) *fileFactory {
	f.prototype.Size = size
	return f
}

func (f *fileFactory) WithSha256(sha256 string) *fileFactory {
	f.prototype.Sha256 = sha256
	return f
}

func (f *fileFactory) WithSha256Enc(sha256Enc string) *fileFactory {
	f.prototype.Sha256Enc = sha256Enc
	return f
}

func (f *fileFactory) WithMediaKey(mediaKey string) *fileFactory {
	f.prototype.MediaKey = mediaKey
	return f
}

func (f *fileFactory) WithPath(path string) *fileFactory {
	f.prototype.Path = path
	return f
}

func (f *fileFactory) WithDirectPath(directPath string) *fileFactory {
	f.prototype.DirectPath = directPath
	return f
}

func (f *fileFactory) WithURL(url string) *fileFactory {
	f.prototype.URL = url
	return f
}

func (f *fileFactory) WithWidth(width *uint32) *fileFactory {
	f.prototype.Width = width
	return f
}

func (f *fileFactory) WithHeight(height *uint32) *fileFactory {
	f.prototype.Height = height
	return f
}

func (f *fileFactory) WithDuration(duration *uint32) *fileFactory {
	f.prototype.Duration = duration
	return f
}

func (f *fileFactory) WithPages(pages *uint32) *fileFactory {
	f.prototype.Pages = pages
	return f
}

func (f *fileFactory) WithThumbnail(thumbnail *file.ImageFile) *fileFactory {
	f.prototype.Thumbnail = thumbnail
	return f
}

func (f *fileFactory) WithCreatedAt(createdAt *time.Time) *fileFactory {
	f.prototype.CreatedAt = createdAt
	return f
}

func (f *fileFactory) WithUpdatedAt(updatedAt *time.Time) *fileFactory {
	f.prototype.UpdatedAt = updatedAt
	return f
}

func (f *fileFactory) WithInstanceID(instanceID *string) *fileFactory {
	f.prototype.InstanceID = instanceID
	return f
}

func (f *fileFactory) Image() *fileFactory {
	f.prototype.Mime = "image/png"
	f.prototype.Extension = "png"
	f.prototype.Size = 2048
	width := uint32(800)
	height := uint32(600)
	f.prototype.Width = &width
	f.prototype.Height = &height
	f.prototype.Duration = nil
	f.prototype.Pages = nil
	return f
}

func (f *fileFactory) Create() *file.File {
	// clona para não vazar referência
	t := *f.prototype

	if t.CreatedAt == nil {
		now := time.Now().UTC()
		t.CreatedAt = &now
	}

	if t.UpdatedAt == nil {
		now := time.Now().UTC()
		t.UpdatedAt = &now
	}

	if t.ID == "" {
		t.ID = uuid.NewString()
	}

	if t.Sha256 == "" {
		t.Sha256 = FakeSha256()
	}

	if t.Sha256Enc == "" {
		t.Sha256Enc = FakeSha256()
	}

	if t.MediaKey == "" {
		t.MediaKey = uuid.NewString()
	}

	if t.Path == "" {
		t.Path = "/path/to/" + t.ID
	}

	if t.DirectPath == "" {
		t.DirectPath = "https://files.example.com/" + t.ID
	}

	if t.URL == "" {
		t.URL = "https://cdn.example.com/" + t.ID
	}

	return &t
}

func (f *fileFactory) CreateMany(n int) []*file.File {
	files := make([]*file.File, 0, n)
	for i := 0; i < n; i++ {
		now := time.Now().Add(time.Duration(i) * time.Minute).UTC()
		f.WithCreatedAt(&now)
		f.WithUpdatedAt(&now)
		files = append(files, f.Create())
	}
	return files
}
