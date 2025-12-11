package led

import (
	"context"
	"machine"

	"github.com/magradze/gonnect/drivers/gpio"
	"github.com/magradze/gonnect/event"
	"github.com/magradze/gonnect/pkg/logger"
	"github.com/magradze/gonnect/registry"
)

// ModuleName is the unique identifier for this component.
const ModuleName = "status_led"

// LedModule controls an LED based on system events.
// It listens for "app/command/toggle" events to change state.
type LedModule struct {
	pin *gpio.Pin
}

// init registers the module with the framework automatically on import.
func init() {
	registry.RegisterModule(&LedModule{})
}

// Init claims the hardware resource (GPIO).
func (l *LedModule) Init() error {
	// We use the onboard LED for demonstration.
	// On ESP32, this is usually GPIO 2.
	// On Pico, it's connected to the internal LED pin.
	p, err := gpio.New(machine.LED, machine.PinOutput, ModuleName)
	if err != nil {
		return err
	}
	l.pin = p
	return nil
}

// Start listens for events and updates the hardware state.
func (l *LedModule) Start(ctx context.Context) {
	// Subscribe to the toggle event
	events := event.Subscribe("app/command/toggle")
	logger.Info("[%s] Listening for toggle events...", ModuleName)

	for {
		select {
		case <-ctx.Done():
			// System is shutting down
			return
		case evt := <-events:
			// React to the event
			logger.Debug("[%s] Received event from %s: %v", ModuleName, evt.Source, evt.Payload)
			l.pin.Toggle()
		}
	}
}

// Stop releases the hardware resource.
func (l *LedModule) Stop() error {
	if l.pin != nil {
		// Turn off before releasing
		l.pin.Low()
		return l.pin.Close()
	}
	return nil
}

// Name returns the module identifier.
func (l *LedModule) Name() string {
	return ModuleName
}