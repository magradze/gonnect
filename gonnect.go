// gonnect.go
package gonnect

import "context"

// Service is a semantic alias for any type registered in the Service Locator.
// It serves as a marker to indicate that a struct is intended to be shared.
type Service = any

// Module represents an independent, lifecycle-managed unit (e.g., Drivers, Sensors).
type Module interface {
	// Init configures the module and claims necessary resources.
	// It is called synchronously on the main thread during boot.
	// Critical Setup: Failure here (returning error) halts the system immediately.
	Init() error

	// Start executes the module's main logic.
	// It is launched in a separate Goroutine by the Engine.
	// Daemons (background tasks) must block here until ctx.Done() is closed.
	// Short-lived tasks may perform work and return immediately.
	Start(ctx context.Context)

	// Stop performs cleanup (e.g., disabling hardware, flushing buffers).
	// It is called synchronously during system shutdown.
	// Note: Avoid long blocking operations here to ensure a quick restart.
	Stop() error

	// Name returns the unique identifier for logging and registration.
	// Ideally, return a string constant to avoid heap allocation.
	Name() string
}