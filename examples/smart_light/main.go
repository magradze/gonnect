package main

import (
	"time"

	"github.com/magradze/gonnect/engine"
	"github.com/magradze/gonnect/pkg/logger"

	// Register modules
	_ "github.com/magradze/gonnect/examples/smart_light/modules/smart_button"
	_ "github.com/magradze/gonnect/examples/smart_light/modules/smart_led"
)

func main() {
	// USB Delay
	time.Sleep(2 * time.Second)

	// Enable logs
	logger.SetLevel(logger.LevelDebug)

	// Start Engine
	app := engine.New(nil)
	app.Run()
}