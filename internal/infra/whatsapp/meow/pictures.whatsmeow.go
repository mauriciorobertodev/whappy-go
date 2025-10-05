package meow

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (g *WhatsmeowGateway) GetPictureURL(ctx context.Context, inst *instance.Instance, phoneOrJID string, preview bool, isCommunity bool) (string, error) {
	l := app.GetWhatsappLogger()

	l.Debug("Getting profile picture URL", "instance", inst.ID, "phoneOrJID", phoneOrJID, "preview", preview, "isCommunity", isCommunity)

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return "", err
	}

	if !strings.Contains(phoneOrJID, "@") {
		phoneOrJID = phoneOrJID + "@" + types.DefaultUserServer
	}

	jid, err := types.ParseJID(phoneOrJID)
	if err != nil {
		return "", err
	}

	l.Debug("Looking for picture", "JID", jid.String())

	ch := make(chan *types.ProfilePictureInfo, 1)
	errCh := make(chan error, 1)

	go func() {
		info, err := client.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
			Preview:     preview,
			ExistingID:  "",
			IsCommunity: isCommunity,
		})
		if err != nil {
			errCh <- err
			return
		}
		ch <- info
	}()

	select {
	case info := <-ch:
		if info == nil || info.URL == "" {
			return "", whatsapp.ErrPictureNotFound
		}

		l.Info("Profile picture URL", "URL", info.URL)

		return info.URL, nil
	case err := <-errCh:
		if errors.Is(err, whatsmeow.ErrProfilePictureUnauthorized) {
			return "", whatsapp.ErrHiddenProfilePicture
		}

		if errors.Is(err, whatsmeow.ErrProfilePictureNotSet) {
			return "", whatsapp.ErrNoProfilePicture
		}

		return "", err
	case <-time.After(5 * time.Second):
		return "", whatsapp.ErrPictureNotFound
	}
}
