// engine/engine.go
package engine

import (
	"context"

	// დავამატეთ მთავარი პაკეტის იმპორტი
	"github.com/magradze/gonnect"
	"github.com/magradze/gonnect/config"
	"github.com/magradze/gonnect/pkg/logger"
	"github.com/magradze/gonnect/registry"
)

// Engine is the central orchestrator of the framework.
// It manages the bootstrap process, module lifecycle, and configuration injection.
type Engine struct {
	Config     *config.Manager
	shutdownCh chan struct{}
}

// New creates a new Engine instance.
// 'store' is optional; pass nil if persistence is not required.
func New(store config.Store) *Engine {
	e := &Engine{
		shutdownCh: make(chan struct{}),
	}
	if store != nil {
		e.Config = config.NewManager(store)
	}
	return e
}

// Run executes the main application loop.
// 1. It initializes all registered modules via Init().
// 2. It starts all modules via Start() in separate goroutines.
// 3. It blocks indefinitely until Shutdown() is called.
// 4. It triggers Stop() on all modules for graceful shutdown.
func (e *Engine) Run() {
	logger.Info("Gonnect Engine starting...")

	modules := registry.GetModules()
	logger.Debug("Found %d registered modules", len(modules))

	// --- Phase 1: Initialization ---
	// Init is synchronous. If any module fails to initialize, the system halts.
	// This ensures we don't start with a broken state (e.g., failed hardware lock).
	for _, m := range modules {
		logger.Debug("Initializing module: %s", m.Name())
		if err := m.Init(); err != nil {
			logger.Error("FATAL: Failed to initialize module '%s': %v", m.Name(), err)
			// In embedded systems, failing Init is usually unrecoverable.
			// Panicking here is the correct behavior to trigger a Watchdog Timer (WDT) reset if configured.
			panic(err)
		}
	}
	logger.Info("All modules initialized successfully")

	// --- Phase 2: Startup ---
	// Create a root context that we can cancel upon shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, m := range modules {
		// Launch each module in its own goroutine.
		// We wrap the call to handle panics and ensure the loop variable 'm' is captured correctly.
		
		// FIX: registry.Module -> gonnect.Module
		go func(mod gonnect.Module) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("CRITICAL: Panic recovered in module '%s': %v", mod.Name(), r)
				}
			}()
			mod.Start(ctx)
		}(m)
	}
	logger.Info("System is running")

	// --- Phase 3: Runtime Loop ---
	// Block here forever.
	// In bare-metal embedded systems, there is no OS signal to wait for.
	// The loop exits only if Shutdown() is explicitly called (e.g., by an OTA update process).
	<-e.shutdownCh
	logger.Warn("Shutdown signal received. Stopping modules...")

	// --- Phase 4: Graceful Shutdown ---
	// Cancel the context to notify all modules to exit their loops.
	cancel()

	// Execute Stop() methods for cleanup.
	for _, m := range modules {
		logger.Debug("Stopping module: %s", m.Name())
		if err := m.Stop(); err != nil {
			logger.Error("Error stopping module '%s': %v", m.Name(), err)
		}
	}

	logger.Info("Gonnect Engine stopped.")
}

// Shutdown triggers a graceful shutdown of the engine.
// It unblocks the Run() loop, cancels the context, and stops all modules.
// Useful for OTA updates, deep sleep preparation, or soft restarts.
func (e *Engine) Shutdown() {
	close(e.shutdownCh)
}