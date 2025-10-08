package input

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	"github.com/mauriciorobertodev/whappy-go/internal/utils"
)

type GetToken struct {
	ID string `json:"id"`
}

func (inp *GetToken) Validate() error {
	if !utils.IsUUID(inp.ID) {
		return token.ErrInvalidID
	}

	return nil
}
