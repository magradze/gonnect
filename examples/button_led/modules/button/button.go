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
	PollRate   = 50 * time.Millisecond
)

type ButtonModule struct {
	pin *gpio.Pin
}

func init() {
	registry.RegisterModule(&ButtonModule{})
}

func (b *ButtonModule) Init() error {
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

	logger.Info("%s Polling button state...", logger.Tag(ModuleName))

	lastState := true

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentState := b.pin.Get()
			if lastState && !currentState {
				logger.Debug("%s Button pressed. Publishing toggle event.", logger.Tag(ModuleName))
				
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