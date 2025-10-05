package token

import (
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

func NewHasher(config *config.TokenConfig) token.TokenHasher {
	if config.IsBcrypt() {
		return NewBcryptHasher()
	}

	return NewSimpleHasher()
}

func NewGenerator() token.TokenGenerator {
	return NewSimpleTokenGenerator()
}
