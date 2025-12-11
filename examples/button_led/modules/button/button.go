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
	PollRate   = 50 * time.Millisecond // Check button every 50ms
)

// ButtonModule monitors a physical button and publishes events on press.
type ButtonModule struct {
	pin *gpio.Pin
}

func init() {
	registry.RegisterModule(&ButtonModule{})
}

func (b *ButtonModule) Init() error {
	// We assume the button is connected to GPIO 0 (Boot button on ESP32)
	// or a specific pin on other boards.
	// Using PinInputPullup ensures we don't need external resistors.
	// Note: Change machine.GPIO0 to your specific board's button pin if needed.
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

	// lastState assumes PullUp: True means NOT pressed, False means PRESSED.
	lastState := true 

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentState := b.pin.Get()

			// Detect falling edge (Pressed)
			// True -> False transition
			if lastState && !currentState {
				logger.Debug("[%s] Button pressed! Publishing event...", ModuleName)
				
				// Publish the event. Payload can be anything, here just a string "click".
				event.Publish("app/command/toggle", "click", ModuleName)
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