package token

import (
	"time"

	"github.com/google/uuid"
)

type TokenGenerator interface {
	Generate(id, instanceID string) (string, error)
}

type TokenHasher interface {
	Hash(input string) (string, error)
	Compare(hash, input string) bool
}

type Token struct {
	ID         string // opcional, UUID público
	raw        string // só usado na criação e no header
	Hash       string // hash ou raw, dependendo da estratégia
	InstanceID string // ID da instância associada ao token
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

func NewToken(generator TokenGenerator, hasher TokenHasher, instanceID string) (*Token, error) {
	id, _ := uuid.NewV7()

	raw, err := generator.Generate(id.String(), instanceID)
	if err != nil {
		return nil, err
	}

	hash, err := hasher.Hash(raw)
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:         id.String(),
		raw:        raw,
		Hash:       hash,
		InstanceID: instanceID,
	}, nil
}

func NewTokenFromHash(id, instanceID, hash string) *Token {
	return &Token{
		ID:         id,
		raw:        "",
		Hash:       hash,
		InstanceID: instanceID,
	}
}

func (t *Token) Matches(input string, hasher TokenHasher) bool {
	return hasher.Compare(t.Hash, input)
}

func (t *Token) FullToken() string {
	return t.ID + "|" + t.raw
}

func (t *Token) GetRaw() string {
	return t.raw
}

func (t *Token) GetHash() string {
	return t.Hash
}

func (t *Token) GetMasked() string {
	if len(t.raw) < 20 {
		return "*****"
	}

	first5 := t.raw[:5]
	last5 := t.raw[len(t.raw)-5:]

	return first5 + "****" + last5
}
