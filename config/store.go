// config/store.go
package config

import "errors"

// ErrNoConfig is returned when the store is empty, uninitialized, or the key is missing.
var ErrNoConfig = errors.New("config: not found")

// Store defines the persistence layer for configuration data.
// It abstracts the underlying storage mechanism (NVS, EEPROM, LittleFS, SD Card).
type Store interface {
	// Load retrieves the raw configuration bytes from persistent storage.
	// It must return ErrNoConfig if no data exists.
	// Note: This allocates memory for the result slice. Since config loading
	// usually happens once at startup, this heap allocation is acceptable.
	Load() ([]byte, error)

	// Save persists the raw configuration bytes to storage.
	// Drivers should implement this atomically (e.g., write-verify or swap)
	// to prevent data corruption if power is lost during the write.
	Save(data []byte) error

	// Clear removes the configuration data (Factory Reset).
	// This should physically erase the data or invalidate the key.
	Clear() error
}