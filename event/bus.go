// event/bus.go
package event

import (
	"sync"
	"time"

	"github.com/magradze/gonnect/pkg/logger"
)

// DefaultBufferSize defines the capacity of the subscription channels.
const DefaultBufferSize = 10

// Event represents a message passed through the bus.
type Event struct {
	Topic     string
	Value     int64
	Payload   interface{}
	Source    string
	Timestamp int64
}

// Bus manages the subscription and publication of events.
type Bus struct {
	mu          sync.Mutex
	subscribers map[string][]chan Event
}

// defaultBus is the global instance.
var defaultBus = &Bus{}

func (b *Bus) ensureInit() {
	if b.subscribers == nil {
		b.subscribers = make(map[string][]chan Event)
	}
}

// Subscribe registers a listener for a specific topic.
func Subscribe(topic string) <-chan Event {
	return defaultBus.Subscribe(topic)
}

// Publish broadcasts an event in a non-blocking manner.
func Publish(topic string, value int64, payload interface{}, source string) int {
	return defaultBus.Publish(topic, value, payload, source)
}

// Subscribe (instance method).
func (b *Bus) Subscribe(topic string) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.ensureInit()

	ch := make(chan Event, DefaultBufferSize)

	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make([]chan Event, 0, 2)
	}

	b.subscribers[topic] = append(b.subscribers[topic], ch)
	logger.Debug("EventBus: New subscriber for '%s'", topic)

	return ch
}

// Publish (instance method) - Non-blocking.
func (b *Bus) Publish(topic string, value int64, payload interface{}, source string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribers == nil {
		return 0
	}

	subscribers, found := b.subscribers[topic]
	if !found || len(subscribers) == 0 {
		return 0
	}

	evt := Event{
		Topic:     topic,
		Value:     value,
		Payload:   payload,
		Source:    source,
		Timestamp: time.Now().UnixNano(),
	}

	dropped := 0

	for _, ch := range subscribers {
		select {
		case ch <- evt:
			// Delivered
		default:
			dropped++
			logger.Warn("EventBus: Dropped '%s' from '%s'", topic, source)
		}
	}

	return dropped
}

// PublishBlocking (instance method).
func (b *Bus) PublishBlocking(topic string, value int64, payload interface{}, source string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribers == nil {
		return
	}

	subscribers, found := b.subscribers[topic]
	if !found || len(subscribers) == 0 {
		return
	}

	evt := Event{
		Topic:     topic,
		Value:     value,
		Payload:   payload,
		Source:    source,
		Timestamp: time.Now().UnixNano(),
	}

	for _, ch := range subscribers {
		ch <- evt
	}
}