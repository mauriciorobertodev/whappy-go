package blocklist

type BlocklistChange struct {
	JID     string `json:"jid"`
	LID     string `json:"lid"`
	Blocked bool   `json:"blocked"`
}
