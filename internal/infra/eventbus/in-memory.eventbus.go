package eventbus

import (
	"sync"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

type InMemoryEventBus struct {
	mu          sync.RWMutex
	handlers    map[events.EventName][]chan events.Event
	allHandlers []chan events.Event
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers:    make(map[events.EventName][]chan events.Event),
		allHandlers: make([]chan events.Event, 0),
	}
}

func (b *InMemoryEventBus) Publish(event events.Event) {
	l := app.GetEventBusLogger()
	l.Debug("Publishing event", "event", event.Name)

	b.mu.RLock()
	defer b.mu.RUnlock()

	if chans, ok := b.handlers[event.Name]; ok {
		for _, ch := range chans {
			select {
			case ch <- event: // envia sem bloquear
			default:
				l.Warn("Event channel full, dropping event", "event", event.Name)
			}
		}
	}

	for _, ch := range b.allHandlers {
		select {
		case ch <- event:
		default:
			l.Warn("Global event channel full, dropping event", "event", event.Name)
		}
	}

	l.Info("Event published", "event", event.Name)
}

func (b *InMemoryEventBus) Subscribe(eventName events.EventName, handler events.EventHandler) {
	ch := make(chan events.Event, 100) // buffer de 100 eventos

	b.mu.Lock()
	b.handlers[eventName] = append(b.handlers[eventName], ch)
	b.mu.Unlock()

	go func() {
		for event := range ch {
			handler(event)
		}
	}()
}

func (b *InMemoryEventBus) SubscribeAll(handler events.EventHandler) {
	ch := make(chan events.Event, 100)

	b.mu.Lock()
	b.allHandlers = append(b.allHandlers, ch)
	b.mu.Unlock()

	go func() {
		for event := range ch {
			handler(event)
		}
	}()
}
