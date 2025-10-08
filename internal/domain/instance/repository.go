package instance

type InstanceQueryOptions struct {
	ID     *string         `db:"id"`
	JID    *string         `db:"jid"`
	LID    *string         `db:"lid"`
	Name   *string         `db:"name"`
	Phone  *string         `db:"phone"`
	Status *InstanceStatus `db:"status"`
	Device *string         `db:"device"`

	Limit   *int   `db:"limit"`
	OrderBy string `db:"order_by"`
	SortBy  string `db:"sort_by"`
}

type InstanceQueryOption func(*InstanceQueryOptions)

type InstanceRepository interface {
	Insert(inst *Instance) error
	InsertMany(insts []*Instance) error

	Update(inst *Instance) error

	Get(opts ...InstanceQueryOption) (*Instance, error)
	List(opts ...InstanceQueryOption) ([]*Instance, error)

	Delete(opts ...InstanceQueryOption) error

	Count(opts ...InstanceQueryOption) int
}

func WhereID(id string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.ID = &id
	}
}

func WhereJID(jid string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.JID = &jid
	}
}

func WhereLID(lid string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.LID = &lid
	}
}

func WhereName(name string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.Name = &name
	}
}

func WherePhone(phone string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.Phone = &phone
	}
}

func WhereStatus(status InstanceStatus) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.Status = &status
	}
}

func WhereDevice(device string) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.Device = &device
	}
}

func Limit(n int) InstanceQueryOption {
	return func(o *InstanceQueryOptions) {
		o.Limit = &n
	}
}
