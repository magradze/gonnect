// examples/button_led/main.go
package main

import (
	"time"

	"github.com/magradze/gonnect/engine"
	"github.com/magradze/gonnect/pkg/logger"

	// Import modules for side-effect registration.
	// The init() functions in these packages will register the modules with the framework.
	// Note: We use the blank identifier "_" because we don't call functions from these packages directly.
	_ "github.com/magradze/gonnect/examples/button_led/modules/button"
	_ "github.com/magradze/gonnect/examples/button_led/modules/led"
)

func main() {
	// Optional: Wait for the USB Serial connection to establish.
	// On boards with native USB (like Pico, ESP32-S2/C3), the console needs
	// a moment to connect, otherwise startup logs are missed.
	time.Sleep(2 * time.Second)

	// Set Log Level to Debug.
	// This ensures we see the event bus traffic and module state changes in the console.
	logger.SetLevel(logger.LevelDebug)

	// Create a new engine instance without persistent storage (nil).
	// For this simple example, we don't need to save/load config.
	app := engine.New(nil)

	// Run the engine.
	// This will:
	// 1. Initialize "status_led" and "user_button" (calls Init()).
	// 2. Start them in separate Goroutines (calls Start()).
	// 3. Block forever (calls select{} internally).
	app.Run()
}