package service

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/objects"
)

type SessionService struct {
	instanceRepo instance.InstanceRepository
	whatsapp     whatsapp.WhatsAppGateway
	eventbus     events.EventBus
}

func NewSessionService(instRepo instance.InstanceRepository, whatsapp whatsapp.WhatsAppGateway, eventbus events.EventBus) *SessionService {
	return &SessionService{
		instanceRepo: instRepo,
		whatsapp:     whatsapp,
		eventbus:     eventbus,
	}
}

func (s *SessionService) Pair(ctx context.Context, inst *instance.Instance) *app.AppError {
	if err := inst.CanPair(); err != nil {
		return app.TranslateError("session service pair", err)
	}

	eventsCh, err := s.whatsapp.PairingQrCode(ctx, inst)
	if err != nil {
		return app.TranslateError("session service pair", err)
	}

	go func() {
		for evt := range eventsCh {
			switch evt.Type {
			case whatsapp.QRCodeGenerated:
				s.eventbus.Publish(inst.AttachQRCode(evt.Code))
			case whatsapp.PairingSuccess:
				if err := inst.CanLoginWith(evt.Phone); err != nil {
					s.eventbus.Publish(inst.FailPairing(instance.FailPairingConflictCode, evt.Phone, err))
					break
				}

				s.eventbus.Publish(inst.LoginWith(evt.Phone, evt.JID, evt.LID, evt.Device))
				s.eventbus.Publish(inst.Connect())
			case whatsapp.PairingTimeout:
				s.eventbus.Publish(inst.FailPairing(instance.FailPairingTimeoutCode, "", evt.Error))
			case whatsapp.PairingError:
				switch evt.Error {
				case whatsapp.ErrClientOutdated:
					s.eventbus.Publish(inst.FailPairing(instance.FailPairingClientOutdatedCode, "", evt.Error))
				case whatsapp.ErrScannedWithoutMultiDevice:
					s.eventbus.Publish(inst.FailPairing(instance.FailPairingWithoutMultideviceCode, "", evt.Error))
				default:
					s.eventbus.Publish(inst.FailPairing(instance.FailPairingUnknownCode, "", evt.Error))
				}
			}
			s.instanceRepo.Update(inst)
		}
	}()

	go s.eventbus.Publish(inst.StartPairing())

	if err := s.instanceRepo.Update(inst); err != nil {
		return app.NewDatabaseError("session service pair", err)
	}

	return nil
}

func (s *SessionService) Connect(ctx context.Context, inst *instance.Instance) *app.AppError {
	if err := inst.CanConnect(); err != nil {
		return app.TranslateError("session service connect", err)
	}

	if err := s.whatsapp.Connect(ctx, inst); err != nil {
		inst.MarkDisconnected() // Not published event of disconnect, because already before was not connected
		go s.eventbus.Publish(inst.EventConnectionFailed(err.Error()))
		return app.TranslateError("session service connect", err)
	}

	if err := s.instanceRepo.Update(inst); err != nil {
		return app.NewDatabaseError("session service connect", err)
	}

	go s.eventbus.Publish(inst.Connect())

	return nil
}

func (s *SessionService) Disconnect(ctx context.Context, inst *instance.Instance) *app.AppError {
	if err := inst.CanDisconnect(); err != nil {
		return app.TranslateError("session service disconnect", err)
	}

	if err := s.whatsapp.Disconnect(ctx, inst); err != nil {
		return app.TranslateError("session service disconnect", err)
	}

	if err := s.instanceRepo.Update(inst); err != nil {
		return app.NewDatabaseError("session service disconnect", err)
	}

	go s.eventbus.Publish(inst.Disconnect())

	return nil
}

func (s *SessionService) Logout(ctx context.Context, inst *instance.Instance) *app.AppError {
	if err := inst.CanLogout(); err != nil {
		return app.TranslateError("session service logout", err)
	}

	if err := s.whatsapp.Logout(ctx, inst); err != nil {
		return app.TranslateError("session service logout", err)
	}

	if err := s.instanceRepo.Update(inst); err != nil {
		return app.NewDatabaseError("session service logout", err)
	}

	go s.eventbus.Publish(inst.Logout())

	return nil
}

func (s *SessionService) Ping(ctx context.Context, inst *instance.Instance) (*objects.Ping, *app.AppError) {
	waPing, err := s.whatsapp.Ping(ctx, inst)
	if err != nil {
		return nil, app.TranslateError("session service ping", err)
	}

	return &objects.Ping{
		Status:          inst.Status,
		IsLoggedIn:      waPing.IsLoggedIn,
		IsConnected:     waPing.IsConnected,
		LastLoginAt:     inst.LastLoginAt,
		LastConnectedAt: inst.LastConnectedAt,
		BannedAt:        inst.BannedAt,
		BanExpiresAt:    inst.BanExpiresAt,
	}, nil
}

func (s *SessionService) QrCode(ctx context.Context, inst *instance.Instance) (string, *app.AppError) {
	if !inst.Status.IsPairing() {
		return "", app.NewAppError("session service qrcode", app.CodeInstanceNotPairing, instance.ErrInstanceNotPairing)
	}

	if inst.Status.IsLoggedIn() {
		return "", app.NewAppError("session service qrcode", app.CodeInstanceAlreadyPaired, instance.ErrInstanceAlreadyPaired)
	}

	if inst.LastQRCode == nil || *inst.LastQRCode == "" {
		return "", app.NewAppError("session service qrcode", app.CodeInstanceNoQrCode, instance.ErrInstanceNoQrCode)
	}

	return *inst.LastQRCode, nil
}
