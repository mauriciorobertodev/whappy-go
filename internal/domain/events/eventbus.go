package events

type EventBus interface {
	Publish(event Event)
	Subscribe(name EventName, handler EventHandler)
	SubscribeAll(handler EventHandler)
}

type EventHandler func(Event)
