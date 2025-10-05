package models

import (
	"time"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
)

type SQLToken struct {
	ID         string    `db:"id"`
	InstanceID string    `db:"instance_id"`
	TokenHash  string    `db:"token_hash"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (s *SQLToken) ToEntity() *token.Token {
	return &token.Token{
		ID:         s.ID,
		InstanceID: s.InstanceID,
		Hash:       s.TokenHash,
		CreatedAt:  s.CreatedAt.UTC(),
		UpdatedAt:  s.UpdatedAt.UTC(),
	}
}

func FromTokenEntity(inst *token.Token) (*SQLToken, error) {
	return &SQLToken{
		ID:         inst.ID,
		InstanceID: inst.InstanceID,
		TokenHash:  inst.GetHash(),
		CreatedAt:  inst.CreatedAt.UTC(),
		UpdatedAt:  inst.UpdatedAt.UTC(),
	}, nil
}
