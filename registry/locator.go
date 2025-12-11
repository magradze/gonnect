// registry/locator.go
package registry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/logger"
)

var (
	// ErrServiceNotFound is returned when the requested service name is missing.
	ErrServiceNotFound = errors.New("registry: service not found")
	// ErrTypeMismatch is returned when the service exists but is not of the requested type.
	ErrTypeMismatch = errors.New("registry: service type mismatch")
)

// defaultLocator is the singleton instance.
var defaultLocator = &locator{}

type locator struct {
	// We use sync.Mutex instead of RWMutex.
	// On single-core MCUs, RWMutex adds binary bloat with no parallel performance benefit.
	// Service lookups are fast map reads, so a simple Mutex is sufficient.
	mu       sync.Mutex
	services map[string]interface{}
}

func (l *locator) ensureInit() {
	if l.services == nil {
		l.services = make(map[string]interface{})
	}
}

// RegisterService adds a service implementation to the registry.
// Returns an error if the name is already taken.
func RegisterService(name string, service interface{}) error {
	defaultLocator.mu.Lock()
	defer defaultLocator.mu.Unlock()

	defaultLocator.ensureInit()

	if _, exists := defaultLocator.services[name]; exists {
		return fmt.Errorf("registry: service '%s' already exists", name)
	}

	defaultLocator.services[name] = service
	logger.Debug("Service registered: '%s'", name)
	return nil
}

// UnregisterService removes a service. Safe to call if not found.
func UnregisterService(name string) {
	defaultLocator.mu.Lock()
	defer defaultLocator.mu.Unlock()

	if defaultLocator.services == nil {
		return
	}

	if _, exists := defaultLocator.services[name]; exists {
		delete(defaultLocator.services, name)
		logger.Debug("Service unregistered: '%s'", name)
	}
}

// GetServiceTyped retrieves a strongly-typed instance of a service.
// T is the expected interface or struct type.
//
// Usage:
//
//	mqtt, err := registry.GetServiceTyped[MQTTClient]("mqtt_main")
func GetServiceTyped[T any](name string) (T, error) {
	defaultLocator.mu.Lock()
	defer defaultLocator.mu.Unlock()

	var zero T

	if defaultLocator.services == nil {
		return zero, ErrServiceNotFound
	}

	raw, exists := defaultLocator.services[name]
	if !exists {
		return zero, ErrServiceNotFound
	}

	// Runtime Type Assertion.
	// Note: In TinyGo, this requires Runtime Type Information (RTTI),
	// which adds a small amount to the binary size, but is unavoidable for a dynamic registry.
	typed, ok := raw.(T)
	if !ok {
		return zero, ErrTypeMismatch
	}

	return typed, nil
}