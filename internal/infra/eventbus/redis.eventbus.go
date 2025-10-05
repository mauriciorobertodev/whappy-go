package eventbus

import (
	"context"
	"encoding/json"

	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
	"github.com/redis/go-redis/v9"
)

type RedisEventBus struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisEventBus(config *config.EventBusConfig) *RedisEventBus {
	return &RedisEventBus{
		client: redis.NewClient(&redis.Options{Addr: config.GetRedisAddress()}),
		ctx:    context.Background(),
	}
}

func (b *RedisEventBus) Publish(event events.Event) {
	data, _ := json.Marshal(event)
	b.client.Publish(b.ctx, string(event.Name), data)
}

func (b *RedisEventBus) Subscribe(name events.EventName, handler events.EventHandler) {
	sub := b.client.Subscribe(b.ctx, string(name))
	go func() {
		for msg := range sub.Channel() {
			var e events.Event
			_ = json.Unmarshal([]byte(msg.Payload), &e)
			handler(e)
		}
	}()
}

func (b *RedisEventBus) SubscribeAll(handler events.EventHandler) {
	pubsub := b.client.PSubscribe(b.ctx, "*")

	go func() {
		defer pubsub.Close()

		ch := pubsub.Channel()
		for msg := range ch {
			var e events.Event
			_ = json.Unmarshal([]byte(msg.Payload), &e)

			handler(e)
		}
	}()
}
