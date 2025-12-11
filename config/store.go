package config

import "errors"

// ErrNoConfig is returned when the store is empty or uninitialized.
var ErrNoConfig = errors.New("configuration not found")

// Store defines the persistence layer for configuration data.
// It abstracts the underlying storage mechanism (NVS, EEPROM, FileSystem).
// Drivers must implement this interface to be used by the Config Manager.
type Store interface {
	// Load retrieves the raw configuration bytes from persistent storage.
	// It returns ErrNoConfig if no data exists or the storage is empty.
	Load() ([]byte, error)

	// Save persists the raw configuration bytes to storage.
	// This operation should be atomic where possible to prevent data corruption.
	Save(data []byte) error
}