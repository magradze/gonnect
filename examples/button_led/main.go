package main

import (
	"github.com/magradze/gonnect/engine"

	// Import modules for side-effect registration.
	// The init() functions in these packages will register the modules with the framework.
	_ "github.com/magradze/gonnect/examples/button_led/modules/button"
	_ "github.com/magradze/gonnect/examples/button_led/modules/led"
)

func main() {
	// Create a new engine instance without persistent storage (nil).
	// For this simple example, we don't need to save/load config.
	app := engine.New(nil)

	// Run the engine.
	// This will:
	// 1. Initialize "status_led" and "user_button".
	// 2. Start them in separate Goroutines.
	// 3. Block forever (or until a termination signal is received).
	app.Run()
}