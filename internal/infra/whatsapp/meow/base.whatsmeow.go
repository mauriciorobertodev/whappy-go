package meow

import (
	"context"
	"errors"
	"sync"

	"github.com/mauriciorobertodev/whappy-go/internal/app/cache"
	"github.com/mauriciorobertodev/whappy-go/internal/app/storage"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/instance"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

var (
	ErrServerStatus420 = errors.New("server returned error 420")
)

type WhatsmeowGateway struct {
	clients   sync.Map // map[string]*whatsmeow.Client
	container *sqlstore.Container
	storage   storage.Storage
	eventbus  events.EventBus
	cache     cache.Cache
}

func New(ctx context.Context, config *config.DatabaseConfig, storage storage.Storage, eventbus events.EventBus, cache cache.Cache) *WhatsmeowGateway {
	wmContainer, err := sqlstore.New(ctx, config.CodeDriver(), config.GetDSN(), nil)
	if err != nil {
		panic("❌ Failed to create whatsapp database:" + err.Error())
	}

	return &WhatsmeowGateway{
		container: wmContainer,
		storage:   storage,
		eventbus:  eventbus,
		cache:     cache,
	}
}

func (g *WhatsmeowGateway) getClient(id string) (*whatsmeow.Client, bool) {
	client, ok := g.clients.Load(id)
	if !ok {
		return nil, false
	}
	return client.(*whatsmeow.Client), true
}

func (g *WhatsmeowGateway) getOnlineClient(instID string) (*whatsmeow.Client, error) {
	val, ok := g.clients.Load(instID)
	if !ok {
		return nil, instance.ErrInstanceNotConnected
	}

	client := val.(*whatsmeow.Client)
	if !client.IsConnected() {
		return nil, instance.ErrInstanceNotConnected
	}

	return client, nil
}
