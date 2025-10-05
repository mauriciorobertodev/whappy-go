package file

import (
	"time"
)

type FileRepository interface {
	Insert(file *File) error
	InsertMany(files []*File) error

	Update(file *File) error

	Get(opts ...FileQueryOption) (*File, error)
	List(opts ...FileQueryOption) ([]*File, error)

	Delete(opts ...FileQueryOption) error

	Count(opts ...FileQueryOption) (uint64, error)
}

type FileQueryOptions struct {
	ID            *string    `db:"id"`
	InstanceID    *string    `db:"instance_id"`
	Cursor        *time.Time `db:"cursor"`
	Sha256        *string    `db:"sha256"`
	Limit         *int       `db:"limit"`
	OrderBy       string     `db:"order_by"`
	SortBy        string     `db:"sort_by"`
	HasThumbnail  *bool      `db:"has_thumbnail"`
	WithThumbnail bool       `db:"with_thumbnail"`
}

type FileQueryOption func(*FileQueryOptions)

func WhereHasThumbnail() FileQueryOption {
	return func(o *FileQueryOptions) {
		has := true
		o.HasThumbnail = &has
	}
}

func WhereInstanceID(id string) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.InstanceID = &id
	}
}

func WithCursor(t *time.Time, limit int) FileQueryOption {
	return func(q *FileQueryOptions) {
		q.Cursor = t
		q.Limit = &limit
	}
}

func WhereSha256(sha256 string) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.Sha256 = &sha256
	}
}

func WhereID(id string) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.ID = &id
	}
}

func WhereDoesNotHaveThumbnail() FileQueryOption {
	return func(o *FileQueryOptions) {
		has := false
		o.HasThumbnail = &has
	}
}

func WithLimit(limit int) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.Limit = &limit
	}
}

func OrderByCreatedAt(desc bool) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.OrderBy = "created_at"
		o.SortBy = "ASC"
		if desc {
			o.SortBy = "DESC"
		}
	}
}

func OrderByUpdatedAt(desc bool) FileQueryOption {
	return func(o *FileQueryOptions) {
		o.OrderBy = "updated_at"
		o.SortBy = "ASC"
		if desc {
			o.SortBy = "DESC"
		}
	}
}

func WithThumbnail() FileQueryOption {
	return func(o *FileQueryOptions) {
		o.WithThumbnail = true
	}
}
