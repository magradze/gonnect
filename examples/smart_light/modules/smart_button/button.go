package smart_button

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
	ModuleName = "smart_button"
	Topic      = "input/command"
	
	// Constants for timing logic
	DebounceTime  = 50 * time.Millisecond
	DoubleGap     = 300 * time.Millisecond // Max time between clicks for double click
	LongPressTime = 800 * time.Millisecond // Time to hold for long press
)

// Command codes (passed as Event.Value)
const (
	CmdSingleClick = 1
	CmdDoubleClick = 2
	CmdLongPress   = 3
)

type SmartButton struct {
	pin *gpio.Pin
}

func init() {
	registry.RegisterModule(&SmartButton{})
}

func (b *SmartButton) Init() error {
	// Boot button (GPIO 0 on ESP32/Pico generally, or check your board)
	p, err := gpio.New(machine.GPIO0, machine.PinInputPullup, ModuleName)
	if err != nil {
		return err
	}
	b.pin = p
	return nil
}

func (b *SmartButton) Start(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Millisecond) // Fast polling for responsiveness
	defer ticker.Stop()

	logger.Info("%s Ready. Try Double Click or Long Press!", logger.Tag(ModuleName))

	// State variables
	var (
		isPressed      bool
		pressStartTime time.Time
		lastReleaseTime time.Time
		clickCount     int
		longPressSent  bool
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Read current hardware state (Low = Pressed because of PullUp)
			currentPressed := !b.pin.Get()

			now := time.Now()

			// --- Logic: Button Just Pressed ---
			if currentPressed && !isPressed {
				isPressed = true
				pressStartTime = now
				longPressSent = false
			}

			// --- Logic: Holding Button (Long Press Check) ---
			if isPressed {
				duration := now.Sub(pressStartTime)
				if duration > LongPressTime && !longPressSent {
					logger.Debug("%s Long Press Detected!", logger.Tag(ModuleName))
					event.Publish(Topic, CmdLongPress, nil, ModuleName)
					longPressSent = true
					// Reset click count to avoid confusion on release
					clickCount = 0 
				}
			}

			// --- Logic: Button Just Released ---
			if !currentPressed && isPressed {
				isPressed = false
				
				// Ignore release if it was a long press
				if !longPressSent {
					clickCount++
					lastReleaseTime = now
				}
			}

			// --- Logic: Click Timeout (Detect Single vs Double) ---
			// If button is released, and we have clicks pending, check time gap
			if !isPressed && clickCount > 0 {
				timeSinceRelease := now.Sub(lastReleaseTime)

				if clickCount == 2 {
					// Double Click Detected immediately
					logger.Debug("%s Double Click!", logger.Tag(ModuleName))
					event.Publish(Topic, CmdDoubleClick, nil, ModuleName)
					clickCount = 0
				} else if timeSinceRelease > DoubleGap {
					// Time expired, so it was just a Single Click
					logger.Debug("%s Single Click", logger.Tag(ModuleName))
					event.Publish(Topic, CmdSingleClick, nil, ModuleName)
					clickCount = 0
				}
			}
		}
	}
}

func (b *SmartButton) Stop() error {
	if b.pin != nil {
		return b.pin.Close()
	}
	return nil
}

func (b *SmartButton) Name() string {
	return ModuleName
}