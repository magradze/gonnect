// resource/types.go
package resource

import "strconv"

// Type represents the classification of a hardware resource.
// We use uint8 to minimize the size of the 'resourceKey' struct in the lock map.
type Type uint8

const (
	// GPIO represents General Purpose Input/Output pins.
	GPIO Type = iota
	// I2C represents Inter-Integrated Circuit buses.
	I2C
	// SPI represents Serial Peripheral Interface buses.
	SPI
	// UART represents Universal Asynchronous Receiver-Transmitter ports.
	UART
	// ADC represents Analog-to-Digital Converter channels.
	ADC
	// PWM represents Pulse Width Modulation channels.
	PWM
	// Timer represents hardware timers.
	Timer
	// DMA represents Direct Memory Access channels.
	DMA
	// limit is used for boundary checking in the String method.
	limit
)

// ID represents the numeric identifier of a resource (e.g., Pin Number, Bus ID).
// uint16 is sufficient for almost all MCUs (65535 pins/ports is plenty),
// and it packs efficiently into memory structures compared to 'int'.
type ID uint16

// typeNames provides a look-up table for string representation.
// This is more efficient than a large switch statement in TinyGo.
var typeNames = [...]string{
	"GPIO",
	"I2C",
	"SPI",
	"UART",
	"ADC",
	"PWM",
	"Timer",
	"DMA",
}

// String returns the string representation of the resource type.
func (t Type) String() string {
	if t < limit {
		return typeNames[t]
	}
	return "Unknown(" + strconv.Itoa(int(t)) + ")"
}