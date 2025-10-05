package consumer

import (
	"os"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mdp/qrterminal"
)

type DevConsumer struct {
}

func NewDevConsumer() *DevConsumer {
	return &DevConsumer{}
}

func (c *DevConsumer) Handler(evt events.Event) {
	if evt.Name == instance.EventPairingQRCode {
		payload := evt.Payload.(instance.PayloadInstanceQRCodeGenerated)
		qrterminal.GenerateHalfBlock(payload.QRCode, qrterminal.L, os.Stdout)
	}
}
