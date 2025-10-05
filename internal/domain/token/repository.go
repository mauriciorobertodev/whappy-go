package token

type TokenRepository interface {
	FindByID(id string) (*Token, error)
	FindByInstanceID(instanceID string) ([]*Token, error)
	Insert(token *Token) error
	Delete(id string) error
}
