package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/blocklist"

type Block struct {
	PhoneOrJID string `json:"phoneOrJID"`
}

func (b *Block) Validate() error {
	if b.PhoneOrJID == "" {
		return blocklist.ErrInvalidJID
	}

	return nil
}

type Unblock struct {
	PhoneOrJID string `json:"phoneOrJID"`
}

func (u *Unblock) Validate() error {
	if u.PhoneOrJID == "" {
		return blocklist.ErrInvalidJID
	}

	return nil
}
