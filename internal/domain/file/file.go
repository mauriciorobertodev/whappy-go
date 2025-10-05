package file

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`

	Size      uint64 `json:"size"`
	Sha256    string `json:"sha256"`
	Sha256Enc string `json:"sha256_enc,omitempty"`
	MediaKey  string `json:"media_key,omitempty"`

	Path       string `json:"path"`
	DirectPath string `json:"direct_path,omitempty"`
	URL        string `json:"url,omitempty"`

	Width    *uint32 `json:"width,omitempty"`
	Height   *uint32 `json:"height,omitempty"`
	Duration *uint32 `json:"duration,omitempty"`
	Pages    *uint32 `json:"pages,omitempty"`

	Thumbnail *ImageFile `json:"thumbnail,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	InstanceID *string `json:"instance_id,omitempty"`
}

type ImageFile struct {
	File
	Width  *uint32 `json:"width"`
	Height *uint32 `json:"height"`
}

func (i *ImageFile) HasDimensions() bool {
	return i.Width != nil && i.Height != nil
}

func (i *ImageFile) IsJPEG() bool {
	return strings.EqualFold(i.Mime, "image/jpeg") || strings.EqualFold(i.Mime, "image/jpg")
}

func (i *ImageFile) IsPNG() bool {
	return strings.EqualFold(i.Mime, "image/png")
}

func (i *ImageFile) ToFile() *File {
	return &i.File
}

type VideoFile struct {
	File
	Width    *uint32 `json:"width"`
	Height   *uint32 `json:"height"`
	Duration *uint32 `json:"duration"`
}

type AudioFile struct {
	File
	Duration *uint32 `json:"duration"`
}

type VoiceFile struct {
	File
	Duration *uint32 `json:"duration"`
}

func NewFile(name, mime, extension, path, sha256 string, size uint64, width, height, duration, pages *uint32) *File {
	id, _ := uuid.NewV7()

	if extension == "" && mime != "" {
		extension = DetectExtension(mime)
	}

	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	return &File{
		ID:        id.String(),
		Name:      name,
		Mime:      mime,
		Size:      size,
		Sha256:    sha256,
		Extension: extension,
		Path:      path,
		Width:     width,
		Height:    height,
		Duration:  duration,
		Pages:     pages,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
}

func NewFromMime(mime string) *File {
	uuid, _ := uuid.NewV7()
	extension := DetectExtension(mime)
	name := fmt.Sprintf("%s.%s", uuid, extension)
	path := fmt.Sprintf("%s.%s", uuid, extension)

	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	return &File{
		ID:        uuid.String(),
		Name:      name,
		Mime:      mime,
		Extension: extension,

		Size:   0,
		Sha256: "",

		Path: path,
		URL:  "",

		Width:    nil,
		Height:   nil,
		Duration: nil,
		Pages:    nil,

		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,

		InstanceID: nil,
	}
}

func (f *File) HasThumbnail() bool {
	return f.Thumbnail != nil
}

func (f *File) HasDimensions() bool {
	return f.Width != nil && f.Height != nil
}

func (f *File) HasDuration() bool {
	return f.Duration != nil
}

func (f *File) HasPages() bool {
	return f.Pages != nil
}

func (f *File) CanBeImage() error {
	if !strings.HasPrefix(f.Mime, "image/") {
		return ErrFileCannotBeImage
	}

	return nil
}

func (f *File) ToImageFile() (*ImageFile, error) {
	if err := f.CanBeImage(); err != nil {
		return nil, err
	}

	return &ImageFile{
		File:   *f,
		Width:  f.Width,
		Height: f.Height,
	}, nil
}

func (f *File) CanBeVideo() error {
	if !strings.HasPrefix(f.Mime, "video/") {
		return ErrFileCannotBeVideo
	}

	return nil
}

func (f *File) ToVideoFile() (*VideoFile, error) {
	if err := f.CanBeVideo(); err != nil {
		return nil, err
	}

	return &VideoFile{
		File:     *f,
		Width:    f.Width,
		Height:   f.Height,
		Duration: f.Duration,
	}, nil
}

func (f *File) CanBeAudio() error {
	if !strings.HasPrefix(f.Mime, "audio/") {
		return ErrFileCannotBeAudio
	}

	return nil
}

func (f *File) ToAudioFile() (*AudioFile, error) {
	if err := f.CanBeAudio(); err != nil {
		return nil, err
	}

	return &AudioFile{
		File:     *f,
		Duration: f.Duration,
	}, nil
}

func (f *File) CanBeVoice() error {
	// In real the voice message needs the opus codec, but for now we will accept these mimes. i really don't want check the codec in the backend
	if f.Mime != "audio/ogg" && f.Mime != "audio/opus" && f.Mime != "audio/ogg; codecs=opus" && f.Mime != "application/ogg" {
		return ErrFileCannotBeVoice
	}

	return nil
}

func (f *File) ToVoiceFile() (*VoiceFile, error) {
	if err := f.CanBeVoice(); err != nil {
		return nil, err
	}

	return &VoiceFile{
		File:     *f,
		Duration: f.Duration,
	}, nil
}

type Metadata struct {
	Name     *string
	Mime     *string
	Width    *uint32
	Height   *uint32
	Duration *uint32
	Pages    *uint32
}

func (f *File) UpdateMeta(m Metadata) {
	if m.Name != nil && *m.Name != "" {
		f.Name = *m.Name
	}

	if m.Mime != nil && *m.Mime != "" {
		f.Mime = *m.Mime
		f.Extension = DetectExtension(*m.Mime)
	}

	f.Width = m.Width
	f.Height = m.Height
	f.Duration = m.Duration
	f.Pages = m.Pages
	updatedAt := time.Now().UTC()
	f.UpdatedAt = &updatedAt
}

func (f *File) SetMime(mime string) {
	f.Mime = mime
	f.Extension = DetectExtension(mime)
	updatedAt := time.Now().UTC()
	f.UpdatedAt = &updatedAt
}
