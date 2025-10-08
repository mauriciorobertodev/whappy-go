package fake

import (
	"sync"

	"github.com/mauriciorobertodev/whappy-go/internal/app"
	"github.com/mauriciorobertodev/whappy-go/internal/domain/events"
)

type FakeEventBus struct {
	mu          sync.RWMutex
	handlers    map[events.EventName][]chan events.Event
	allHandlers []chan events.Event
	published   []events.Event
}

func NewFakeEventBus() *FakeEventBus {
	return &FakeEventBus{
		handlers:    make(map[events.EventName][]chan events.Event),
		allHandlers: make([]chan events.Event, 0),
	}
}

func (b *FakeEventBus) Publish(event events.Event) {
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

	b.published = append(b.published, event)

	l.Info("Event published", "event", event.Name)
}

func (b *FakeEventBus) Subscribe(eventName events.EventName, handler events.EventHandler) {
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

func (b *FakeEventBus) SubscribeAll(handler events.EventHandler) {
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

func (b *FakeEventBus) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = make(map[events.EventName][]chan events.Event)
	b.allHandlers = []chan events.Event{}
}

// FAKE

func (b *FakeEventBus) Published() []events.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.published
}

func (b *FakeEventBus) ClearPublished() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.published = []events.Event{}
}

func (b *FakeEventBus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = make(map[events.EventName][]chan events.Event)
	b.allHandlers = []chan events.Event{}
	b.published = []events.Event{}
}

func (b *FakeEventBus) HasPublished(name events.EventName) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, e := range b.published {
		if e.Name == name {
			return true
		}
	}
	return false
}
