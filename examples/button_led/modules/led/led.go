// examples/button_led/modules/led/led.go
package led

import (
	"context"
	"machine"

	"github.com/magradze/gonnect/drivers/gpio"
	"github.com/magradze/gonnect/event"
	"github.com/magradze/gonnect/pkg/logger"
	"github.com/magradze/gonnect/registry"
)

const ModuleName = "status_led"

type LedModule struct {
	pin *gpio.Pin
}

func init() {
	registry.RegisterModule(&LedModule{})
}

func (l *LedModule) Init() error {
	targetPin := machine.GPIO13 
	p, err := gpio.New(targetPin, machine.PinOutput, ModuleName)
	if err != nil {
		return err
	}
	l.pin = p
	return nil
}

func (l *LedModule) Start(ctx context.Context) {
	events := event.Subscribe("app/command/toggle")
	
	// ცვლილება: ვიყენებთ logger.Tag()-ს ფერისთვის
	logger.Info("%s Listening for toggle events...", logger.Tag(ModuleName))

	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-events:
			logger.Debug("%s Toggle signal received (Source: %s)", 
				logger.Tag(ModuleName), evt.Source)
			
			l.pin.Toggle()
		}
	}
}

func (l *LedModule) Stop() error {
	if l.pin != nil {
		l.pin.Low()
		return l.pin.Close()
	}
	return nil
}

func (l *LedModule) Name() string {
	return ModuleName
}