package eventbroker

type EventType int

const (
	EventCreatePlayer EventType = iota
)

type Observer interface {
	Observe(eventType EventType, event interface{})
}

type EventBroker interface {
	Broadcast(eventType EventType, event interface{})
	AddObserver(observer Observer, eventTypes []EventType)
}

type EventBrokerImpl struct {
	observations map[EventType][]Observer
}

func NewEventBroker() *EventBrokerImpl {
	return &EventBrokerImpl{
		observations: map[EventType][]Observer{},
	}
}

func (e *EventBrokerImpl) Broadcast(eventType EventType, event interface{}) {
	for _, observer := range e.observations[eventType] {
		observer.Observe(eventType, event)
	}
}

func (e *EventBrokerImpl) AddObserver(observer Observer, eventTypes []EventType) {
	for _, eventType := range eventTypes {
		e.observations[eventType] = append(e.observations[eventType], observer)
	}
}
