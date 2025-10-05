package token

type PayloadTokenRenewed struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Masked bool   `json:"masked"`
}
