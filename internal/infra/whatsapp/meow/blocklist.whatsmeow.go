package meow

import (
	"context"
	"strings"

	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (g *WhatsmeowGateway) GetBlockList(ctx context.Context, inst *instance.Instance) ([]string, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	list, err := client.GetBlocklist()
	if err != nil {
		return nil, err
	}

	blocked := make([]string, 0, len(list.JIDs))
	for _, jid := range list.JIDs {
		blocked = append(blocked, jid.String())
	}

	return blocked, nil
}

func (g *WhatsmeowGateway) UpdateBlockList(ctx context.Context, inst *instance.Instance, phoneOrJID string, action whatsapp.BlocklistAction) ([]string, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(phoneOrJID, "@") {
		phoneOrJID = phoneOrJID + "@" + types.DefaultUserServer
	}

	jid, err := types.ParseJID(phoneOrJID)
	if err != nil {
		return nil, err
	}

	list, err := client.UpdateBlocklist(jid, events.BlocklistChangeAction(action))
	if err != nil {
		return nil, err
	}

	blocked := make([]string, 0, len(list.JIDs))
	for _, jid := range list.JIDs {
		blocked = append(blocked, jid.String())
	}

	return blocked, nil
}
