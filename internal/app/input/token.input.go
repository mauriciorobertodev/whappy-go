package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/token"

type GetToken struct {
	ID string `json:"id"`
}

func (inp *GetToken) Validate() error {
	if inp.ID == "" {
		return token.ErrInvalidID
	}

	return nil
}
