// gonnect.go
package gonnect

import "context"

// Service marks a type as a discoverable service within the framework.
// Any struct that needs to be accessible via the Service Locator must implement this.
type Service interface{}

// Module represents an independent unit of functionality (e.g., WiFi, MQTT, Sensors).
// Every component that needs lifecycle management by the Engine must implement this interface.
type Module interface {
	// Init initializes the module.
	// It is called synchronously during the system bootstrap phase.
	// Returning an error here halts the system startup.
	Init() error

	// Start executes the main logic of the module.
	// It is launched in a separate goroutine by the Engine.
	// The module must respect ctx.Done() for graceful shutdown.
	Start(ctx context.Context)

	// Stop performs necessary cleanup operations (e.g., closing connections, freeing resources).
	// It is called during the system shutdown process.
	Stop() error

	// Name returns the unique identifier of the module.
	// This identifier is used for logging, debugging, and registration.
	Name() string
}