package meow

import (
	"go.mau.fi/whatsmeow/types"
)

func IsLID(jid types.JID) bool {
	return jid.Server == "lid"
}

func IsPN(jid types.JID) bool {
	return jid.Server == "s.whatsapp.net"
}

func IsGroup(jid types.JID) bool {
	return jid.Server == types.GroupServer
}

func IsBroadcast(jid types.JID) bool {
	return jid.Server == types.BroadcastServer
}

func IsNewsletter(jid types.JID) bool {
	return jid.Server == types.NewsletterServer
}

func IsStatus(jid types.JID) bool {
	return jid.String() == types.StatusBroadcastJID.String()
}

func IsUser(jid types.JID) bool {
	return jid.Server == types.DefaultUserServer || jid.Server == types.HiddenUserServer
}

func GetJIDAndLID(JID types.JID) (string, string) {
	jid := types.NewJID(JID.User, JID.Server).String()
	lid := ""

	if IsLID(JID) {
		lid = types.NewJID(JID.User, JID.Server).String()
		jid = ""
	}

	return jid, lid
}
