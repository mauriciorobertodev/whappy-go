package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type BlocklistService struct {
	whatsapp whatsapp.WhatsAppGateway
	eventbus events.EventBus
}

func NewBlocklistService(whatsapp whatsapp.WhatsAppGateway, eventbus events.EventBus) *BlocklistService {
	return &BlocklistService{
		whatsapp: whatsapp,
		eventbus: eventbus,
	}
}

func (s *BlocklistService) GetBlocklist(ctx context.Context, inst *instance.Instance) ([]string, *app.AppError) {
	l := app.GetBlocklistServiceLogger()
	l.Debug("Retrieving blocklist", "instance", inst.ID)

	blocklist, err := s.whatsapp.GetBlockList(ctx, inst)
	if err != nil {
		return []string{}, app.TranslateError("blocklist service", err)
	}

	return blocklist, nil
}

func (s *BlocklistService) Block(ctx context.Context, inst *instance.Instance, inp input.Block) ([]string, *app.AppError) {
	l := app.GetBlocklistServiceLogger()
	l.Debug("Blocking phone or JID", "instance", inst.ID, "phoneOrJID", inp.PhoneOrJID)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("blocklist service", err)
	}

	blocklist, err := s.whatsapp.UpdateBlockList(ctx, inst, inp.PhoneOrJID, whatsapp.BlocklistActionBlock)
	if err != nil {
		return nil, app.TranslateError("blocklist service", err)
	}

	return blocklist, nil
}

func (s *BlocklistService) Unblock(ctx context.Context, inst *instance.Instance, inp input.Unblock) ([]string, *app.AppError) {
	l := app.GetBlocklistServiceLogger()
	l.Debug("Unblocking phone or JID", "instance", inst.ID, "phoneOrJID", inp.PhoneOrJID)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("blocklist service", err)
	}

	blocklist, err := s.whatsapp.UpdateBlockList(ctx, inst, inp.PhoneOrJID, whatsapp.BlocklistActionUnblock)
	if err != nil {
		return nil, app.TranslateError("blocklist service", err)
	}

	return blocklist, nil
}
