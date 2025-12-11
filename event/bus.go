package event

import (
	"sync"

	"github.com/magradze/gonnect/pkg/logger"
)

// DefaultBufferSize defines the capacity of the channel returned by Subscribe.
// This buffer prevents the publisher from blocking if the consumer is slightly slow.
const DefaultBufferSize = 10

// Event represents a message passed through the bus.
type Event struct {
	// Topic is the channel name (e.g., "wifi/status", "sensor/temp").
	Topic string
	// Payload is the actual data (can be anything).
	Payload interface{}
	// Source identifies the module that published the event.
	Source string
}

// Bus manages the subscription and publication of events.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Event
}

// defaultBus is the global instance of the event bus.
var defaultBus = &Bus{
	subscribers: make(map[string][]chan Event),
}

// Subscribe registers a listener for a specific topic.
// It returns a read-only channel that will receive events for that topic.
// The returned channel has a buffer size of DefaultBufferSize.
func Subscribe(topic string) <-chan Event {
	return defaultBus.Subscribe(topic)
}

// Publish broadcasts an event to all subscribers of the given topic.
// It uses a non-blocking send; if a subscriber's channel is full, the event is dropped for that subscriber
// to prevent blocking the entire system.
func Publish(topic string, payload interface{}, source string) {
	defaultBus.Publish(topic, payload, source)
}

// Subscribe (instance method) registers a new channel for the topic.
func (b *Bus) Subscribe(topic string) <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, DefaultBufferSize)
	
	// Create slice if it doesn't exist
	if _, ok := b.subscribers[topic]; !ok {
		b.subscribers[topic] = make([]chan Event, 0)
	}
	
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	logger.Debug("New subscriber registered for topic: '%s'", topic)
	
	return ch
}

// Publish (instance method) sends data to all subscribers.
func (b *Bus) Publish(topic string, payload interface{}, source string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subscribers, found := b.subscribers[topic]
	if !found || len(subscribers) == 0 {
		// No subscribers, just ignore
		return
	}

	evt := Event{
		Topic:   topic,
		Payload: payload,
		Source:  source,
	}

	for _, ch := range subscribers {
		select {
		case ch <- evt:
			// Success
		default:
			// Channel is full. We drop the message to avoid blocking the publisher.
			// In a real-time system, latest data is usually more important than old buffered data.
			logger.Warn("Event Bus: Subscriber channel full for topic '%s'. Dropping event from '%s'.", topic, source)
		}
	}
}