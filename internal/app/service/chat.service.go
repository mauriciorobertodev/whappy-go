package service

import (
	"context"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/chat"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
)

type ChatService struct {
	whatsapp whatsapp.WhatsAppGateway
}

func NewChatService(whatsapp whatsapp.WhatsAppGateway) *ChatService {
	return &ChatService{
		whatsapp,
	}
}

func (s *ChatService) SendPresence(ctx context.Context, inst *instance.Instance, inp input.SendChatPresenceInput) error {
	l := app.GetChatServiceLogger()
	l.Debug("Sending chat presence", "instance", inst.ID, "to", inp.To, "type", inp.Type, "duration", inp.Duration, "wait", inp.Wait)

	if err := inp.Validate(); err != nil {
		return app.TranslateError("chat service", err)
	}

	if err := s.whatsapp.SendChatPresence(ctx, inst, chat.Presence{
		To:   inp.To,
		Type: inp.Type,
	}); err != nil {
		l.Error("Error sending presence", "error", err)
		return app.TranslateError("chat service", err)
	}

	sendPause := func() {
		if err := s.whatsapp.SendChatPresence(ctx, inst, chat.Presence{
			To:   inp.To,
			Type: chat.ChatPresencePaused,
		}); err != nil {
			l.Error("Error sending pause presence", "error", err)
		}
	}

	if inp.Duration == nil || *inp.Duration <= 0 {
		l.Debug("Presence sent with no duration", "instance", inst.ID)
		return nil
	}

	// Normaliza duração
	duration := *inp.Duration
	if duration > whatsapp.MaxPresenceDuration {
		duration = whatsapp.MaxPresenceDuration
	}

	if inp.Wait != nil && *inp.Wait {
		select {
		case <-time.After(time.Duration(duration) * time.Millisecond):
			sendPause()
		case <-ctx.Done():
			l.Warn("Presence wait cancelled by context", "error", ctx.Err())
			return ctx.Err()
		}
	} else {
		go func() {
			select {
			case <-time.After(time.Duration(duration) * time.Millisecond):
				sendPause()
			case <-ctx.Done():
				l.Debug("Presence goroutine cancelled by context")
				return
			}
		}()
	}

	l.Debug("Presence handling completed", "instance", inst.ID)
	return nil
}
