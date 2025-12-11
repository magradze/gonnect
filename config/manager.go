// config/manager.go
package config

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/cbor"
	"github.com/magradze/gonnect/pkg/logger"
)

// Manager handles the lifecycle of configuration data.
// It acts as a bridge between the raw storage driver and the application's struct data.
type Manager struct {
	mu    sync.RWMutex
	store Store
}

// NewManager creates a new instance of the Config Manager.
// 'store' must be a valid implementation of the Store interface (e.g., NVS, File).
func NewManager(store Store) *Manager {
	return &Manager{
		store: store,
	}
}

// Load reads data from storage and unmarshals it into the provided pointer 'v'.
// 'v' must be a pointer to a struct.
// Returns ErrNoConfig if the storage is empty.
func (m *Manager) Load(v interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.store == nil {
		return fmt.Errorf("config manager: no storage driver configured")
	}

	// 1. Read raw bytes
	data, err := m.store.Load()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return ErrNoConfig
	}

	// 2. Decode CBOR
	if err := cbor.Unmarshal(data, v); err != nil {
		logger.Error("Config Manager: failed to decode configuration: %v", err)
		return fmt.Errorf("config decode failed: %w", err)
	}

	logger.Debug("Configuration loaded successfully (%d bytes)", len(data))
	return nil
}

// Save marshals the provided value 'v' into CBOR and writes it to storage.
// This operation is thread-safe.
func (m *Manager) Save(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return fmt.Errorf("config manager: no storage driver configured")
	}

	// 1. Encode to CBOR
	data, err := cbor.Marshal(v)
	if err != nil {
		return fmt.Errorf("config encode failed: %w", err)
	}

	// 2. Write to storage
	if err := m.store.Save(data); err != nil {
		logger.Error("Config Manager: failed to write to storage: %v", err)
		return err
	}

	logger.Info("Configuration saved successfully (%d bytes)", len(data))
	return nil
}