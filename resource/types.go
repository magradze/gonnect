// resource/types.go
package resource

// Type represents the classification of a hardware resource.
// It is used to segregate locks for different hardware subsystems
// to prevent conflicts (e.g., ensuring two drivers don't claim the same GPIO).
type Type int

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
)

// String returns the string representation of the resource type.
// Useful for logging and debugging resource conflicts.
func (t Type) String() string {
	switch t {
	case GPIO:
		return "GPIO"
	case I2C:
		return "I2C"
	case SPI:
		return "SPI"
	case UART:
		return "UART"
	case ADC:
		return "ADC"
	case PWM:
		return "PWM"
	case Timer:
		return "Timer"
	case DMA:
		return "DMA"
	default:
		return "Unknown"
	}
}

// ID represents the numeric identifier of a specific resource within a Type.
// For GPIO, this corresponds to the pin number.
// For buses (I2C/SPI), this corresponds to the port number (0, 1, etc.).
type ID int