package webhook

type SortBy string

const (
	SortByAsc  SortBy = "ASC"
	SortByDesc SortBy = "DESC"
)

type WebhookQueryOptions struct {
	ID         *string `db:"id"`
	Active     *bool   `db:"active"`
	URL        *string `db:"url"`
	InstanceID *string `db:"instance_id"`

	Limit   *int   `db:"limit"`
	OrderBy string `db:"order_by"`
	SortBy  SortBy `db:"sort_by"`
}

type WebhookQueryOption func(*WebhookQueryOptions)

type WebhookRepository interface {
	Insert(webhook *Webhook) error
	InsertMany(webhooks []*Webhook) error

	Update(webhook *Webhook) error

	Get(opts ...WebhookQueryOption) (*Webhook, error)
	List(opts ...WebhookQueryOption) ([]*Webhook, error)

	Delete(opts ...WebhookQueryOption) error
}

func WhereID(id string) WebhookQueryOption {
	return func(o *WebhookQueryOptions) {
		o.ID = &id
	}
}

func WhereActive(active bool) WebhookQueryOption {
	return func(o *WebhookQueryOptions) {
		o.Active = &active
	}
}

func WhereInstanceID(instanceID string) WebhookQueryOption {
	return func(o *WebhookQueryOptions) {
		o.InstanceID = &instanceID
	}
}

func Limit(limit int) WebhookQueryOption {
	return func(o *WebhookQueryOptions) {
		o.Limit = &limit
	}
}

func OrderBy(orderBy string, sortBy SortBy) WebhookQueryOption {
	return func(o *WebhookQueryOptions) {
		o.OrderBy = orderBy
		o.SortBy = sortBy
	}
}
