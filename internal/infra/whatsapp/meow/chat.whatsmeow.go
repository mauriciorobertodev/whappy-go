package meow

import (
	"context"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/chat"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow/types"
)

func (g *WhatsmeowGateway) SendChatPresence(ctx context.Context, inst *instance.Instance, presence chat.Presence) error {
	l := app.GetWhatsappLogger()
	l.Info("Sending chat presence", "to", presence.To, "type", presence.Type)

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return err
	}

	var presenceType types.ChatPresence
	var presenceMedia types.ChatPresenceMedia

	switch presence.Type {
	case chat.ChatPresenceTyping:
		presenceType = types.ChatPresenceComposing
		presenceMedia = types.ChatPresenceMediaText
	case chat.ChatPresenceRecording:
		presenceType = types.ChatPresenceComposing
		presenceMedia = types.ChatPresenceMediaAudio
	default:
		presenceType = types.ChatPresencePaused
	}

	to, err := types.ParseJID(presence.To)
	if err != nil {
		return err
	}

	if err := client.SendChatPresence(to, presenceType, presenceMedia); err != nil {
		return err
	}

	return nil
}
