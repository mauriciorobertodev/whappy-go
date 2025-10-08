package input

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type CreateInstance struct {
	Name string `json:"name"`
}

func (i *CreateInstance) Validate() error {
	if i.Name == "" {
		return instance.ErrNameTooShort
	}

	if len(i.Name) > instance.MaxNameLength {
		return instance.ErrNameTooLong
	}

	return nil
}

type GetInstance struct {
	ID string `json:"id"`
}

func (i *GetInstance) Validate() error {
	if i.ID == "" {
		return instance.ErrInvalidID
	}

	return nil
}

type RenewInstanceToken struct {
	ID string `json:"id"`
}

func (i *RenewInstanceToken) Validate() error {
	if i.ID == "" {
		return instance.ErrInvalidID
	}

	return nil
}
