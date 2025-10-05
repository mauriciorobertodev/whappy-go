package config

type TokenHasher string

const (
	HasherBcrypt TokenHasher = "bcrypt"
	HasherSimple TokenHasher = "simple"
)

func (h TokenHasher) IsValid() bool {
	switch h {
	case HasherBcrypt, HasherSimple:
		return true
	default:
		return false
	}
}

type TokenConfig struct {
	Hasher TokenHasher
}

func LoadTokenConfig() *TokenConfig {
	cfg := &TokenConfig{
		Hasher: TokenHasher(GetEnvString("TOKEN_HASHER", "simple")),
	}

	if !cfg.Hasher.IsValid() {
		panic("Invalid Token Hasher: " + string(cfg.Hasher))
	}

	return cfg
}

func (c *TokenConfig) IsBcrypt() bool {
	return c.Hasher == HasherBcrypt
}

func (c *TokenConfig) IsSimple() bool {
	return c.Hasher == HasherSimple
}
