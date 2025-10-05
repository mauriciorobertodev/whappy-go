package privacy

type PayloadPrivacyChanged struct {
	LastSeen     *PrivacySetting `json:"last_seen,omitempty"`
	Status       *PrivacySetting `json:"status,omitempty"`
	Profile      *PrivacySetting `json:"profile,omitempty"`
	GroupAdd     *PrivacySetting `json:"group_add,omitempty"`
	ReadReceipts *PrivacySetting `json:"read_receipts,omitempty"`
	CallAdd      *PrivacySetting `json:"call_add,omitempty"`
	Online       *PrivacySetting `json:"online,omitempty"`
}
