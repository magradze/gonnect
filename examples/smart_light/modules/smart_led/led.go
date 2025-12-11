package smart_led

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
	ModuleName = "smart_led"
	Topic      = "input/command"
)

// Modes
const (
	ModeOff       = 0
	ModeStrobe    = 1 // Fast blink (Alarm) - from Long Press
	ModeHeartbeat = 2 // Pulse effect - from Double Click
)

type SmartLed struct {
	pin  *gpio.Pin
	mode int
}

func init() {
	registry.RegisterModule(&SmartLed{})
}

func (l *SmartLed) Init() error {
    targetPin := machine.GPIO13 

    p, err := gpio.New(targetPin, machine.PinOutput, ModuleName)
    if err != nil {
        return err
    }
    l.pin = p
    return nil
}

func (l *SmartLed) Start(ctx context.Context) {
	events := event.Subscribe(Topic)
	
	// Base ticker for animation frames (50ms resolution)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	logger.Info("%s Listening...", logger.Tag(ModuleName))

	// Animation counter
	tickCount := 0

	for {
		select {
		case <-ctx.Done():
			return

		// --- Event Handling ---
		case evt := <-events:
			cmd := evt.Value
			switch cmd {
			case 1: // Single Click -> OFF
				logger.Debug("%s Mode: OFF", logger.Tag(ModuleName))
				l.mode = ModeOff
				l.pin.Low()
			case 2: // Double Click -> Heartbeat
				logger.Debug("%s Mode: HEARTBEAT (Pulsing)", logger.Tag(ModuleName))
				l.mode = ModeHeartbeat
			case 3: // Long Press -> Strobe
				logger.Debug("%s Mode: STROBE (Fast Blink)", logger.Tag(ModuleName))
				l.mode = ModeStrobe
			}
			// Reset animation counter on mode change
			tickCount = 0

		// --- Animation Loop ---
		case <-ticker.C:
			tickCount++
			
			switch l.mode {
			case ModeOff:
				// Ensure it stays off
				// (Optional: optimize by checking current state)
				l.pin.Low()

			case ModeStrobe:
				// Fast blink: Toggle every 100ms (2 ticks)
				if tickCount%2 == 0 {
					l.pin.Toggle()
				}

			case ModeHeartbeat:
				// Simulate Heartbeat:
				// Tick 0-2: ON
				// Tick 2-4: OFF
				// Tick 4-6: ON
				// Tick 6-25: OFF (Long pause)
				cycle := tickCount % 25
				if (cycle >= 0 && cycle < 2) || (cycle >= 4 && cycle < 6) {
					l.pin.High()
				} else {
					l.pin.Low()
				}
			}
		}
	}
}

func (l *SmartLed) Stop() error {
	if l.pin != nil {
		l.pin.Low()
		return l.pin.Close()
	}
	return nil
}

func (l *SmartLed) Name() string {
	return ModuleName
}