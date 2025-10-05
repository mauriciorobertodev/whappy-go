package meow

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/file"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow"
)

func (g *WhatsmeowGateway) UploadFile(ctx context.Context, inst *instance.Instance, stream io.ReadCloser, kind whatsapp.MediaKind, mime string) (*file.File, error) {
	l := app.GetWhatsappLogger()

	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	if stream == nil {
		return nil, fmt.Errorf("file reader is nil")
	}

	defer stream.Close()

	up, err := client.UploadReader(ctx, stream, nil, mediaKindToWhatsmeowMediaType(kind))
	if err != nil {
		return nil, err
	}

	l.Info("File uploaded to WhatsApp servers", "url", up.URL)

	return &file.File{
		URL:        up.URL,
		DirectPath: up.DirectPath,
		Size:       up.FileLength,
		MediaKey:   hex.EncodeToString(up.MediaKey),
		Sha256:     hex.EncodeToString(up.FileSHA256),
		Sha256Enc:  hex.EncodeToString(up.FileEncSHA256),
	}, nil
}

func mediaKindToWhatsmeowMediaType(kind whatsapp.MediaKind) whatsmeow.MediaType {
	switch kind {
	case whatsapp.MediaImage:
		return whatsmeow.MediaImage
	case whatsapp.MediaVideo:
		return whatsmeow.MediaVideo
	case whatsapp.MediaAudio:
		return whatsmeow.MediaAudio
	default:
		return whatsmeow.MediaDocument
	}
}
