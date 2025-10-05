package meow

import (
	"context"
	"fmt"
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (g *WhatsmeowGateway) PairingQrCode(ctx context.Context, inst *instance.Instance) (<-chan whatsapp.PairingEvent, error) {
	out := make(chan whatsapp.PairingEvent, 1)

	deviceStore := g.container.NewDevice()
	client := whatsmeow.NewClient(deviceStore, nil)

	ctxQR, cancel := context.WithTimeout(ctx, 240*time.Second)
	qrChan, err := client.GetQRChannel(ctxQR)
	if err != nil {
		cancel()
		return nil, err
	}

	if err := client.Connect(); err != nil {
		cancel()
		return nil, err
	}

	go func() {
		defer func() {
			cancel()
			close(out)
		}()

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				out <- whatsapp.PairingEvent{
					Type: whatsapp.QRCodeGenerated,
					Code: evt.Code,
				}
			case "success":
				device := client.Store.ID
				out <- whatsapp.PairingEvent{
					Type:   whatsapp.PairingSuccess,
					Phone:  device.User,
					JID:    types.NewJID(device.User, device.Server).String(),
					LID:    types.NewJID(client.Store.LID.User, client.Store.LID.Server).String(),
					Device: device.String(),
				}
			case "timeout":
				out <- whatsapp.PairingEvent{
					Type:  whatsapp.PairingTimeout,
					Error: whatsapp.ErrTimeout,
				}
			case "err-client-outdated":
				out <- whatsapp.PairingEvent{
					Type:  whatsapp.PairingError,
					Error: whatsapp.ErrClientOutdated,
				}
			case "err-scanned-without-multidevice":
				out <- whatsapp.PairingEvent{
					Type:  whatsapp.PairingError,
					Error: whatsapp.ErrScannedWithoutMultiDevice,
				}
			default:
				err := fmt.Errorf("unknown error: %s", evt.Event)
				if evt.Error != nil {
					err = fmt.Errorf("%w: %v", err, evt.Error)
				}
				out <- whatsapp.PairingEvent{
					Type:  whatsapp.PairingError,
					Error: err,
				}
			}
		}
	}()

	return out, nil
}

func (g *WhatsmeowGateway) Logout(ctx context.Context, inst *instance.Instance) error {
	client, ok := g.getClient(inst.ID)
	if !ok {
		return instance.ErrInstanceNotConnected
	}

	if client.IsLoggedIn() {
		if err := client.Logout(ctx); err != nil {
			return err
		}
	}
	return nil
}
