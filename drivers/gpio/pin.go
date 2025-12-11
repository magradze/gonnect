// drivers/gpio/pin.go
package gpio

import (
	"machine"

	"github.com/magradze/gonnect/resource"
)

// Pin is a wrapper around the standard machine.Pin.
// It integrates with the Gonnect Resource Manager to ensure exclusive access.
type Pin struct {
	hw    machine.Pin
	id    resource.ID
	owner string
}

// New claims a GPIO pin and configures it.
// It returns an error if the pin is already in use by another module.
//
// Usage:
//
//	led, err := gpio.New(machine.LED, machine.PinOutput, "status_led")
func New(pin machine.Pin, mode machine.PinMode, owner string) (*Pin, error) {
	// Convert machine.Pin (which is usually an int/uint type) to our resource ID
	resID := resource.ID(pin)

	// 1. Attempt to lock the resource via the Manager
	if err := resource.Lock(resource.GPIO, resID, owner); err != nil {
		return nil, err
	}

	// 2. Configure the hardware
	pin.Configure(machine.PinConfig{Mode: mode})

	return &Pin{
		hw:    pin,
		id:    resID,
		owner: owner,
	}, nil
}

// Set sets the pin logic level (High/Low).
func (p *Pin) Set(high bool) {
	p.hw.Set(high)
}

// High sets the pin output to logic high.
func (p *Pin) High() {
	p.hw.High()
}

// Low sets the pin output to logic low.
func (p *Pin) Low() {
	p.hw.Low()
}

// Get reads the current logic level of the pin.
func (p *Pin) Get() bool {
	return p.hw.Get()
}

// Toggle inverts the current pin state.
// Note: This assumes the pin is configured as output.
func (p *Pin) Toggle() {
	if p.Get() {
		p.Low()
	} else {
		p.High()
	}
}

// Close releases the pin resource.
// It allows other modules to claim this pin later.
func (p *Pin) Close() error {
	return resource.Unlock(resource.GPIO, p.id, p.owner)
}