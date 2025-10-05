package file

import "time"

type PayloadFileUploaded struct {
	ID string `json:"id"`

	Path string `json:"path"`
	URL  string `json:"url"`

	Name     *string `json:"name,omitempty"`
	Mime     *string `json:"mime,omitempty"`
	Width    *uint32 `json:"width,omitempty"`
	Height   *uint32 `json:"height,omitempty"`
	Duration *uint32 `json:"duration,omitempty"`
	Pages    *uint32 `json:"pages,omitempty"`

	Size   uint64 `json:"size,omitempty"`
	Sha256 string `json:"sha256,omitempty"`

	UploadedAt time.Time `json:"uploaded_at"`
}

type PayloadFileDeleted struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

type PayloadFileUpdated struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Name      *string   `json:"name,omitempty"`
	Mime      *string   `json:"mime,omitempty"`
	Width     *uint32   `json:"width,omitempty"`
	Height    *uint32   `json:"height,omitempty"`
	Duration  *uint32   `json:"duration,omitempty"`
	Pages     *uint32   `json:"pages,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
