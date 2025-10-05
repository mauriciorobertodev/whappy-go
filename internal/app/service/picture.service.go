package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type PictureService struct {
	whatsapp whatsapp.WhatsAppGateway
}

func NewPictureService(whatsapp whatsapp.WhatsAppGateway) *PictureService {
	return &PictureService{
		whatsapp: whatsapp,
	}
}

func (s *PictureService) Get(ctx context.Context, inst *instance.Instance, inp input.GetPictureInput) (string, *app.AppError) {
	l := app.GetPictureServiceLogger()
	l.Debug("Getting contact", "instance", inst.ID, "phoneOrJID", inp.PhoneOrJID)

	if err := inp.Validate(); err != nil {
		return "", app.TranslateError("picture service", err)
	}

	if inp.Preview == nil {
		inp.Preview = new(bool)
		*inp.Preview = true
	}

	if inp.IsCommunity == nil {
		inp.IsCommunity = new(bool)
		*inp.IsCommunity = false
	}

	pictureURL, err := s.whatsapp.GetPictureURL(ctx, inst, inp.PhoneOrJID, *inp.Preview, *inp.IsCommunity)
	if err != nil {
		return "", app.TranslateError("picture service", err)
	}

	return pictureURL, nil
}
