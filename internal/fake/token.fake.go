package fake

import (
	"time"

	"github.com/google/uuid"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
)

type tokenFactory struct {
	prototype *token.Token
}

func TokenFactory() *tokenFactory {
	return &tokenFactory{
		prototype: &token.Token{
			ID:         "",
			Hash:       "",
			InstanceID: "",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		},
	}
}

// Métodos fluentes
func (f *tokenFactory) WithID(id string) *tokenFactory {
	f.prototype.ID = id
	return f
}

func (f *tokenFactory) WithInstanceID(instanceID string) *tokenFactory {
	f.prototype.InstanceID = instanceID
	return f
}

func (f *tokenFactory) WithHash(hash string) *tokenFactory {
	f.prototype.Hash = hash
	return f
}

// Exemplo de "estado" igual Laravel
func (f *tokenFactory) Expired() *tokenFactory {
	f.prototype.UpdatedAt = time.Now().Add(-24 * time.Hour)
	return f
}

// Criação final
func (f *tokenFactory) Create() *token.Token {
	// clona para não vazar referência
	t := *f.prototype

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}

	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = time.Now().UTC()
	}

	if t.ID == "" {
		t.ID = uuid.NewString()
	}

	if t.Hash == "" {
		t.Hash = uuid.NewString()
	}

	if t.InstanceID == "" {
		t.InstanceID = uuid.NewString()
	}

	return &t
}

func (f *tokenFactory) CreateMany(n int) []*token.Token {
	tokens := make([]*token.Token, n)
	for i := 0; i < n; i++ {
		tokens[i] = f.Create()
	}
	return tokens
}
