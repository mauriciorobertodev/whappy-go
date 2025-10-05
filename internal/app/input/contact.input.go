package input

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/contact"
)

type CheckPhones struct {
	Phones []string `json:"phones"`
}

func (c *CheckPhones) Validate() error {
	if len(c.Phones) == 0 {
		return contact.ErrInvalidJID
	}

	return nil
}

type GetContact struct {
	PhoneOrJID string `json:"phoneOrJID"`
}

func (g *GetContact) Validate() error {
	if g.PhoneOrJID == "" {
		return contact.ErrInvalidJID
	}

	return nil
}
