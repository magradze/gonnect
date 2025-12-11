// engine/engine.go
package engine

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/magradze/gonnect/config"
	"github.com/magradze/gonnect/pkg/logger"
	"github.com/magradze/gonnect/registry"
)

// Engine is the central orchestrator of the framework.
// It manages the bootstrap process, module lifecycle, and configuration injection.
type Engine struct {
	Config *config.Manager
}

// New creates a new Engine instance.
// 'store' is optional; pass nil if persistence is not required.
func New(store config.Store) *Engine {
	e := &Engine{}
	if store != nil {
		e.Config = config.NewManager(store)
	}
	return e
}

// Run executes the main application loop.
// 1. It initializes all registered modules via Init().
// 2. It starts all modules via Start() in separate goroutines.
// 3. It blocks until a termination signal (SIGINT/SIGTERM) is received.
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
			// In embedded systems, we might want to restart or enter safe mode here.
			// For now, we panic to stop execution.
			panic(err)
		}
	}
	logger.Info("All modules initialized successfully")

	// --- Phase 2: Startup ---
	// Create a root context that we can cancel upon shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, m := range modules {
		// We launch each module in its own goroutine.
		// The module is responsible for keeping its loop alive until ctx.Done().
		go m.Start(ctx)
	}
	logger.Info("System is running")

	// --- Phase 3: Runtime Loop & Shutdown Handling ---
	// Block here until we receive an interrupt signal (Ctrl+C or system kill).
	// On some microcontrollers, this might never happen, effectively blocking forever.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for signal
	sig := <-sigChan
	logger.Warn("Received termination signal: %v. Shutting down...", sig)

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

	logger.Info("Gonnect Engine stopped. Goodbye.")
}