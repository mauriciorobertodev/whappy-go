package meow

import (
	"context"
	"strings"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/whatsapp"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/contact"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"go.mau.fi/whatsmeow/types"
)

// Retrieves all contacts from the instance's contact list.
func (g *WhatsmeowGateway) GetContacts(ctx context.Context, inst *instance.Instance) ([]*contact.Contact, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*contact.Contact, 0, len(contacts))

	for JID, c := range contacts {
		result = append(result, contactInfoToContact(inst, c, JID))
	}

	return result, nil
}

// Retrieves a specific contact by JID from the instance's contact list.
func (g *WhatsmeowGateway) GetContact(ctx context.Context, inst *instance.Instance, phoneOrJID string) (*contact.Contact, error) {
	l := app.GetWhatsappLogger()
	l.Info("Getting contact", "instance", inst.ID, "phoneOrJID", phoneOrJID)

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

	l.Info("Looking for contact", "JID", jid.String())

	contactInfo, err := client.Store.Contacts.GetContact(ctx, jid)
	if err != nil {
		return nil, err
	}

	if !contactInfo.Found {
		return nil, contact.ErrNotFound
	}

	return contactInfoToContact(inst, contactInfo, jid), nil
}

func (g *WhatsmeowGateway) CheckPhones(ctx context.Context, inst *instance.Instance, phones []string) ([]whatsapp.PhoneStatus, error) {
	client, err := g.getOnlineClient(inst.ID)
	if err != nil {
		return nil, err
	}

	checkedPhones := make([]whatsapp.PhoneStatus, 0, len(phones))

	checked, err := client.IsOnWhatsApp(phones)
	if err != nil {
		return nil, err
	}

	for index, res := range checked {
		jid := res.JID

		isBusiness := false
		if res.VerifiedName != nil {
			isBusiness = true
		}

		verifiedName := new(string)
		if res.VerifiedName != nil && res.VerifiedName.Details != nil && res.VerifiedName.Details.VerifiedName != nil {
			*verifiedName = *res.VerifiedName.Details.VerifiedName
		}

		checkedPhones = append(checkedPhones, whatsapp.PhoneStatus{
			Original:     phones[index],
			JID:          jid.String(),
			Phone:        jid.User,
			Exists:       res.IsIn,
			IsBusiness:   isBusiness,
			VerifiedName: verifiedName,
		})
	}

	return checkedPhones, nil
}

func contactInfoToContact(inst *instance.Instance, info types.ContactInfo, JID types.JID) *contact.Contact {
	jid, lid := GetJIDAndLID(JID)

	phone := ""
	if jid != "" {
		phone = JID.User
	}

	if info.RedactedPhone != "" {
		phone = info.RedactedPhone
	}

	return &contact.Contact{
		JID:          jid,
		LID:          lid,
		Phone:        phone,
		FirstName:    info.FirstName,
		FullName:     info.FullName,
		PushName:     info.PushName,
		BusinessName: info.BusinessName,
		IsBusiness:   info.BusinessName != "",
		IsMe:         JID.String() == inst.JID || JID.String() == inst.LID || phone == inst.Phone,
		IsHidden:     lid != "",
	}
}
