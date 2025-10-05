package token

import (
	"crypto/rand"
	"encoding/hex"
)

type SimpleTokenGenerator struct{}

func NewSimpleTokenGenerator() *SimpleTokenGenerator {
	return &SimpleTokenGenerator{}
}

func (g *SimpleTokenGenerator) Generate(tokenId, instanceId string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
