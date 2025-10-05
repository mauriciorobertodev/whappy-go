package events

type PayloadMessageQueued struct {
	ID         string `json:"id"`
	To         string `json:"to"`
	InstanceID string `json:"instance_id"`
}
