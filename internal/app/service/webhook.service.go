package service

import (
	"context"
	"errors"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/webhook"
)

type WebhookService struct {
	webRepo     webhook.WebhookRepository
	bus         events.EventBus
	maxWebhooks int
}

func NewWebhookService(webRepo webhook.WebhookRepository, bus events.EventBus, maxWebhooks int) *WebhookService {
	return &WebhookService{
		webRepo:     webRepo,
		bus:         bus,
		maxWebhooks: maxWebhooks,
	}
}

func (s *WebhookService) CreateWebhook(ctx context.Context, inst *instance.Instance, inp input.CreateWebhook) (*webhook.Webhook, string, *app.AppError) {
	l := app.GetUploadServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, "", app.TranslateError("webhook service", err)
	}

	count, err := s.webRepo.Count(webhook.WhereInstanceID(inst.ID))
	if err != nil {
		return nil, "", app.NewDatabaseError("webhook service", err)
	}

	if s.maxWebhooks > 0 && int(count) >= s.maxWebhooks {
		return nil, "", app.NewAppError("webhook service", app.CodeWebhookMaxWebhooksReached, webhook.ErrMaxWebhooksReached)
	}

	web := webhook.New(inp.URL, inp.Events, inp.Active)
	web.AttachToInstance(inst.ID)

	if err := s.webRepo.Insert(web); err != nil {
		return nil, "", app.TranslateError("webhook service", err)
	}

	l.Info("webhook created", "webhook", web.ID, "instance", inst.ID)

	return web, web.GetSecret(), nil
}

func (s *WebhookService) GetWebhook(ctx context.Context, inst *instance.Instance, inp input.GetWebhook) (*webhook.Webhook, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Debug("getting webhook for instance", "instance", inst.ID, "webhook", inp.ID)

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("webhook service", err)
	}

	web, err := s.webRepo.Get(webhook.WhereInstanceID(inst.ID), webhook.WhereID(inp.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return nil, app.TranslateError("webhook service", err)
		}
		return nil, app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhook retrieved", "instance", inst.ID, "webhook", web.ID)

	return web, nil
}

func (s *WebhookService) GetWebhooks(ctx context.Context, inst *instance.Instance) ([]*webhook.Webhook, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Debug("getting webhooks for instance", "instance", inst.ID)

	webhooks, err := s.webRepo.List(webhook.WhereInstanceID(inst.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return nil, app.TranslateError("webhook service", err)
		}
		return nil, app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhooks retrieved", "instance", inst.ID)

	return webhooks, nil
}

func (s *WebhookService) UpdateWebhook(ctx context.Context, inst *instance.Instance, inp input.UpdateWebhook) (*webhook.Webhook, *app.AppError) {
	l := app.GetUploadServiceLogger()

	l.Debug("validating update webhook input")

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("webhook service", err)
	}

	web, err := s.webRepo.Get(webhook.WhereInstanceID(inst.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return nil, app.TranslateError("webhook service", err)
		}
		return nil, app.NewDatabaseError("webhook service", err)
	}

	web.Update(inp.URL, inp.Events)
	if inp.Active {
		web.Activate()
	} else {
		web.Deactivate()
	}

	if err := s.webRepo.Update(web); err != nil {
		return nil, app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhook updated", "webhook", web.ID, "instance", inst.ID)

	return web, nil
}

func (s *WebhookService) ToggleWebhook(ctx context.Context, inst *instance.Instance, inp input.ToggleWebhook) *app.AppError {
	l := app.GetUploadServiceLogger()

	web, err := s.webRepo.Get(webhook.WhereInstanceID(inst.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return app.TranslateError("webhook service", err)
		}
		return app.NewDatabaseError("webhook service", err)
	}

	if inp.Active {
		web.Activate()
	} else {
		web.Deactivate()
	}

	if err := s.webRepo.Update(web); err != nil {
		return app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhook toggled", "webhook", web.ID, "instance", inst.ID, "active", web.Active)

	return nil
}

func (s *WebhookService) RenewWebhookSecret(ctx context.Context, inst *instance.Instance, inp *input.RenewWebhookSecret) (*webhook.Webhook, string, *app.AppError) {
	l := app.GetUploadServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, "", app.TranslateError("webhook service", err)
	}

	web, err := s.webRepo.Get(webhook.WhereInstanceID(inst.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return nil, "", app.TranslateError("webhook service", err)
		}
		return nil, "", app.NewDatabaseError("webhook service", err)
	}

	web.RenewSecret()

	if err := s.webRepo.Update(web); err != nil {
		return nil, "", app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhook secret renewed", "webhook", web.ID, "instance", inst.ID)

	return web, web.GetSecret(), nil
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, inst *instance.Instance, inp input.DeleteWebhook) *app.AppError {
	l := app.GetUploadServiceLogger()

	if err := inp.Validate(); err != nil {
		return app.TranslateError("webhook service", err)
	}

	err := s.webRepo.Delete(webhook.WhereInstanceID(inst.ID), webhook.WhereID(inp.ID))
	if err != nil {
		if errors.Is(err, webhook.ErrNotFound) {
			return app.TranslateError("webhook service", err)
		}
		return app.NewDatabaseError("webhook service", err)
	}

	l.Info("webhook deleted", "webhook", inp.ID, "instance", inst.ID)

	return nil
}
