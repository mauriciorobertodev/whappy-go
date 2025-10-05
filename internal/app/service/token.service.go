package service

import (
	"context"
	"strings"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/input"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/token"
	c "github.com/mauriciorobertodev/whappy-go/internal/infra/cache"
)

type TokenService struct {
	tokenRepo   token.TokenRepository
	tokenHasher token.TokenHasher
	tokenGen    token.TokenGenerator
	eventbus    events.EventBus
	cache       cache.Cache
}

func NewTokenService(repo token.TokenRepository, hasher token.TokenHasher, gen token.TokenGenerator, eventbus events.EventBus, cache cache.Cache) *TokenService {
	return &TokenService{
		tokenRepo:   repo,
		tokenHasher: hasher,
		tokenGen:    gen,
		eventbus:    eventbus,
		cache:       cache,
	}
}

func (s *TokenService) Get(ctx context.Context, inp input.GetToken) (*token.Token, *app.AppError) {
	l := app.GetTokenServiceLogger()

	if err := inp.Validate(); err != nil {
		return nil, app.TranslateError("token service", err)
	}

	cacheKey := cache.CacheKeyTokenPrefix + inp.ID

	var tok *token.Token

	cachedToken, err := c.Get[token.Token](s.cache, cacheKey)
	if err == nil {
		l.Info("Token found in cache", "id", inp.ID)
		return &cachedToken, nil
	}

	tok, err = s.tokenRepo.FindByID(inp.ID)
	if err != nil {
		return nil, app.NewDatabaseError("token service", err)
	}

	if tok == nil {
		l.Warn("No token found with ID: %s\n", inp.ID)
		return nil, app.TranslateError("token service", token.ErrNotFound)
	}

	go c.Set(s.cache, cacheKey, *tok, cache.DefaultTTL)

	return tok, nil
}

func (s *TokenService) CreateToken(instID string) (*token.Token, *app.AppError) {
	tok, err := token.NewToken(s.tokenGen, s.tokenHasher, instID)
	if err != nil {
		return nil, app.TranslateError("token service", err)
	}

	if err := s.tokenRepo.Insert(tok); err != nil {
		return nil, app.NewDatabaseError("token service", err)
	}

	return tok, nil
}

func (s *TokenService) ValidateFullToken(ctx context.Context, fullToken string) (*token.Token, *app.AppError) {
	l := app.GetTokenServiceLogger()

	if !strings.Contains(fullToken, "|") {
		l.Warn("Token does not contain expected delimiter")
		return nil, app.TranslateError("token service", token.ErrInvalidToken)
	}

	parts := strings.Split(fullToken, "|")
	if len(parts) != 2 {
		l.Warn("Token does not split into two parts as expected")
		return nil, app.TranslateError("token service", token.ErrInvalidToken)
	}

	id := parts[0]
	raw := parts[1]

	tok, appErr := s.Get(ctx, input.GetToken{ID: id})
	if appErr != nil {
		return nil, app.WrapLoc("token service", appErr)
	}

	if !tok.Matches(raw, s.tokenHasher) {
		l.Warn("Token hash does not match for ID", "id", id)
		return nil, app.TranslateError("token service", token.ErrInvalidToken)
	}

	return tok, nil
}

func (s *TokenService) DeleteToken(id string) *app.AppError {

	if err := s.tokenRepo.Delete(id); err != nil {
		return app.NewDatabaseError("token service", err)
	}

	cacheKey := cache.CacheKeyTokenPrefix + id

	go s.cache.Delete(cacheKey)

	return nil
}

func (s *TokenService) RenewToken(ctx context.Context, id string, masking bool) (*token.Token, *app.AppError) {
	tok, appErr := s.Get(ctx, input.GetToken{ID: id})
	if appErr != nil {
		return nil, app.WrapLoc("token service", appErr)
	}

	appErr = s.DeleteToken(tok.ID)
	if appErr != nil {
		return nil, app.WrapLoc("token service", appErr)
	}

	t, appErr := s.CreateToken(tok.InstanceID)
	if appErr != nil {
		return nil, app.WrapLoc("token service", appErr)
	}

	go s.eventbus.Publish(t.EventRenewed(masking))
	// Return the token with the new raw value for the user to see it once
	return t, nil
}

func (s *TokenService) GetByInstanceID(instID string) ([]*token.Token, *app.AppError) {
	tokens, err := s.tokenRepo.FindByInstanceID(instID)
	if err != nil {
		return nil, app.NewDatabaseError("token service", err)
	}

	return tokens, nil
}
