// drivers/gpio/pin.go
package gpio

import (
	"machine"

	"github.com/magradze/gonnect/resource"
)

// Pin is a secure wrapper around the standard machine.Pin.
// It enforces resource locking via the Manager to prevent hardware conflicts.
type Pin struct {
	hw    machine.Pin
	id    resource.ID
	owner string
}

// New claims a GPIO pin, locks it, and configures the hardware mode.
// It returns an error if the pin is already locked by another module.
//
// Usage:
//
//	led, err := gpio.New(machine.LED, machine.PinOutput, "status_led")
func New(pin machine.Pin, mode machine.PinMode, owner string) (*Pin, error) {
	// Cast machine.Pin to our internal uint16 ID type.
	// This works across all TinyGo supported architectures (AVR, ARM, RISC-V).
	resID := resource.ID(uint16(pin))

	// 1. Acquire Lock
	if err := resource.Lock(resource.GPIO, resID, owner); err != nil {
		return nil, err
	}

	// 2. Configure Hardware
	// TinyGo's Configure panics on invalid configuration, which is acceptable at startup.
	pin.Configure(machine.PinConfig{Mode: mode})

	return &Pin{
		hw:    pin,
		id:    resID,
		owner: owner,
	}, nil
}

// Set changes the pin logic level.
func (p *Pin) Set(high bool) {
	p.hw.Set(high)
}

// High sets the pin to logic high (VCC).
func (p *Pin) High() {
	p.hw.High()
}

// Low sets the pin to logic low (GND).
func (p *Pin) Low() {
	p.hw.Low()
}

// Get reads the current logic level.
// Note: Depending on the architecture and PinMode, this reads the Input Register.
func (p *Pin) Get() bool {
	return p.hw.Get()
}

// Toggle inverts the current pin state.
// This is a software-based toggle (Read-Modify-Write).
func (p *Pin) Toggle() {
	p.hw.Set(!p.hw.Get())
}

// Close releases the resource lock.
// The pin hardware state remains unchanged (it does not automatically reset to input).
func (p *Pin) Close() error {
	return resource.Unlock(resource.GPIO, p.id, p.owner)
}