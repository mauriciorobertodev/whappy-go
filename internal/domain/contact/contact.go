package contact

type Contact struct {
	JID          string `json:"jid"`           // The JID (Jabber ID) of the contact
	LID          string `json:"lid"`           // The LID (Local ID) of the contact, if available
	Phone        string `json:"phone"`         // The phone number of the contact
	FirstName    string `json:"first_name"`    // The first name of the contact
	FullName     string `json:"full_name"`     // The full name of the contact
	PushName     string `json:"push_name"`     // The name set by the user
	BusinessName string `json:"business_name"` // The business name of the contact, if it's a business account
	IsBusiness   bool   `json:"is_business"`   // Whether the contact is a business account
	IsMe         bool   `json:"is_me"`         // Whether the contact is the user themselves
	IsHidden     bool   `json:"is_hidden"`     // Whether the contact is a hidden user
	// Picture      string `json:"picture"`       // URL to the contact's avatar image
}
