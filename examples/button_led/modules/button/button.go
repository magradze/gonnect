// examples/button_led/modules/button/button.go
package button

import (
	"context"
	"machine"
	"time"

	"github.com/magradze/gonnect/drivers/gpio"
	"github.com/magradze/gonnect/event"
	"github.com/magradze/gonnect/pkg/logger"
	"github.com/magradze/gonnect/registry"
)

const (
	ModuleName = "user_button"
	// PollRate acts as a software debounce.
	PollRate = 50 * time.Millisecond
)

// ButtonModule monitors a physical button and publishes events on press.
type ButtonModule struct {
	pin *gpio.Pin
}

func init() {
	registry.RegisterModule(&ButtonModule{})
}

func (b *ButtonModule) Init() error {
	// Target pin: GPIO 0 (Common "Boot" button on ESP32 boards).
	// On RP2040/Pico, you might need to change this to a specific pin.
	targetPin := machine.GPIO0

	p, err := gpio.New(targetPin, machine.PinInputPullup, ModuleName)
	if err != nil {
		return err
	}
	b.pin = p
	return nil
}

func (b *ButtonModule) Start(ctx context.Context) {
	ticker := time.NewTicker(PollRate)
	defer ticker.Stop()

	logger.Info("[%s] Polling button state...", ModuleName)

	// lastState assumes PullUp logic:
	// True  (High) = Released
	// False (Low)  = Pressed
	lastState := true

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentState := b.pin.Get()

			// Detect Falling Edge (Transition from High -> Low)
			if lastState && !currentState {
				logger.Debug("[%s] Button pressed. Publishing toggle event.", ModuleName)

				// Publish the event using the new Zero-Allocation signature.
				// Topic: "app/command/toggle"
				// Value: 1 (int64 representation of a "trigger" or "true")
				// Payload: nil (No complex data needed, saving Heap allocation)
				// Source: ModuleName
				event.Publish("app/command/toggle", 1, nil, ModuleName)
			}

			lastState = currentState
		}
	}
}

func (b *ButtonModule) Stop() error {
	if b.pin != nil {
		return b.pin.Close()
	}
	return nil
}

func (b *ButtonModule) Name() string {
	return ModuleName
}