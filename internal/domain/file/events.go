package file

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

// To listen all events, use the prefix "file:*"
const (
	EventFileUploaded events.EventName = "file:uploaded"
	EventFileDeleted  events.EventName = "file:deleted"
	EventFileUpdated  events.EventName = "file:updated"
)

func (f *File) EventUploaded(instanceID *string) events.Event {
	return events.New(
		EventFileUploaded,
		PayloadFileUploaded{
			ID: f.ID,

			Path: f.Path,
			URL:  f.URL,

			Name:     &f.Name,
			Mime:     &f.Mime,
			Width:    f.Width,
			Height:   f.Height,
			Duration: f.Duration,
			Pages:    f.Pages,

			Size:   f.Size,
			Sha256: f.Sha256,

			UploadedAt: *f.CreatedAt,
		},
		instanceID,
	)
}

func (f *File) EventDeleted(instanceID *string) events.Event {
	return events.New(
		EventFileDeleted,
		PayloadFileDeleted{
			ID:        f.ID,
			DeletedAt: time.Now(),
		},
		instanceID,
	)
}

func (f *File) EventUpdated(instanceID *string) events.Event {
	return events.New(
		EventFileUpdated,
		PayloadFileUpdated{
			ID:        f.ID,
			URL:       f.URL,
			Name:      &f.Name,
			Mime:      &f.Mime,
			Width:     f.Width,
			Height:    f.Height,
			Duration:  f.Duration,
			Pages:     f.Pages,
			UpdatedAt: *f.UpdatedAt,
		},
		f.InstanceID,
	)
}
