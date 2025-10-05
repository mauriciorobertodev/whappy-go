package input

import "github.com/mauriciorobertodev/whappy-go/internal/domain/picture"

type GetPictureInput struct {
	PhoneOrJID  string `json:"phone_or_jid"`
	Preview     *bool  `json:"preview"`
	IsCommunity *bool  `json:"is_community"`
}

func (inp *GetPictureInput) Validate() error {
	if inp.PhoneOrJID == "" {
		return picture.ErrInvalidJID
	}

	return nil
}
