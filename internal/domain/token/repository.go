package token

type TokenQueryOptions struct {
	ID         *string `db:"id"`
	InstanceID *string `db:"instance_id"`

	Limit   *int   `db:"limit"`
	OrderBy string `db:"order_by"`
	SortBy  string `db:"sort_by"`
}

type TokenQueryOption func(*TokenQueryOptions)

type TokenRepository interface {
	FindByID(id string) (*Token, error)
	FindByInstanceID(instanceID string) ([]*Token, error)
	Insert(token *Token) error
	Delete(id string) error
	Count(opts ...TokenQueryOption) int
}

func WhereID(id string) TokenQueryOption {
	return func(o *TokenQueryOptions) {
		o.ID = &id
	}
}

func WhereInstanceID(instanceID string) TokenQueryOption {
	return func(o *TokenQueryOptions) {
		o.InstanceID = &instanceID
	}
}

func Limit(limit int) TokenQueryOption {
	return func(o *TokenQueryOptions) {
		o.Limit = &limit
	}
}
