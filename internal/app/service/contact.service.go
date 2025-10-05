package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/contact"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type ContactService struct {
	whatsapp whatsapp.WhatsAppGateway
}

func NewContactService(whatsapp whatsapp.WhatsAppGateway) *ContactService {
	return &ContactService{
		whatsapp: whatsapp,
	}
}

func (s *ContactService) Check(ctx context.Context, inst *instance.Instance, inp input.CheckPhones) ([]whatsapp.PhoneStatus, *app.AppError) {
	l := app.GetContactServiceLogger()
	l.Debug("Checking phones", "instance", inst.ID, "phones", inp.Phones)

	if err := inp.Validate(); err != nil {
		return []whatsapp.PhoneStatus{}, app.TranslateError("contact service", err)
	}

	checked, err := s.whatsapp.CheckPhones(ctx, inst, inp.Phones)
	if err != nil {
		return []whatsapp.PhoneStatus{}, app.TranslateError("contact service", err)
	}

	return checked, nil
}

func (s *ContactService) GetContact(ctx context.Context, inst *instance.Instance, inp input.GetContact) (*contact.Contact, *app.AppError) {
	l := app.GetContactServiceLogger()
	l.Debug("Getting contact", "instance", inst.ID, "phoneOrJID", inp.PhoneOrJID)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("contact service", err)
	}

	c, err := s.whatsapp.GetContact(ctx, inst, inp.PhoneOrJID)
	if err != nil {
		return nil, app.TranslateError("contact service", err)
	}

	return c, nil
}

func (s *ContactService) GetContacts(ctx context.Context, inst *instance.Instance) ([]*contact.Contact, *app.AppError) {
	l := app.GetContactServiceLogger()
	l.Debug("Getting contacts", "instance", inst.ID)

	contacts, err := s.whatsapp.GetContacts(ctx, inst)
	if err != nil {
		return nil, app.TranslateError("contact service", err)
	}

	return contacts, nil
}
