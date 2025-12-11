package registry

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/logger"
)

// defaultLocator is the singleton instance of the service registry.
var defaultLocator = &locator{
	services: make(map[string]interface{}),
}

// locator manages the registration and discovery of services.
// It allows modules to communicate via direct API calls without hard dependencies.
// We use RWMutex because services are read frequently (Get) but registered rarely (Register).
type locator struct {
	mu       sync.RWMutex
	services map[string]interface{}
}

// RegisterService adds a service implementation to the registry.
// If a service with the same name already exists, it returns an error.
// The 'service' argument can be any interface or struct pointer.
func RegisterService(name string, service interface{}) error {
	defaultLocator.mu.Lock()
	defer defaultLocator.mu.Unlock()

	if _, exists := defaultLocator.services[name]; exists {
		return fmt.Errorf("service registry: service '%s' is already registered", name)
	}

	defaultLocator.services[name] = service
	logger.Debug("Service registered: '%s'", name)
	return nil
}

// UnregisterService removes a service from the registry.
// Safe to call even if the service does not exist.
func UnregisterService(name string) {
	defaultLocator.mu.Lock()
	defer defaultLocator.mu.Unlock()

	if _, exists := defaultLocator.services[name]; exists {
		delete(defaultLocator.services, name)
		logger.Debug("Service unregistered: '%s'", name)
	}
}

// GetService retrieves a raw interface for the requested service name.
// It returns (nil, false) if the service is not found.
// Note: Prefer using GetServiceTyped for type safety.
func GetService(name string) (interface{}, bool) {
	defaultLocator.mu.RLock()
	defer defaultLocator.mu.RUnlock()

	svc, ok := defaultLocator.services[name]
	return svc, ok
}

// GetServiceTyped retrieves a strongly-typed instance of a service.
// T is the expected interface or struct type.
//
// Usage:
//
//	mqtt, err := registry.GetServiceTyped[MQTTClient]("mqtt_main")
//
// Returns an error if the service is not found or if the type assertion fails.
func GetServiceTyped[T any](name string) (T, error) {
	defaultLocator.mu.RLock()
	defer defaultLocator.mu.RUnlock()

	var zero T // Zero value for T (e.g., nil for pointers)

	raw, exists := defaultLocator.services[name]
	if !exists {
		return zero, fmt.Errorf("service registry: service '%s' not found", name)
	}

	// Dynamic Type Assertion
	typed, ok := raw.(T)
	if !ok {
		return zero, fmt.Errorf("service registry: service '%s' exists but does not match expected type", name)
	}

	return typed, nil
}