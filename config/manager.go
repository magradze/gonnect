// config/manager.go
package config

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/cbor"
	"github.com/magradze/gonnect/pkg/logger"
)

// Validator is an optional interface for configuration structs.
type Validator interface {
	Validate() error
}

// Manager handles the lifecycle of persistent configuration data.
type Manager struct {
	// RWMutex -> Mutex
	mu    sync.Mutex
	store Store
}

func NewManager(store Store) *Manager {
	return &Manager{
		store: store,
	}
}

// Load reads data from storage.
func (m *Manager) Load(v interface{}) error {
	m.mu.Lock() // RLock -> Lock
	defer m.mu.Unlock() // RUnlock -> Unlock

	if m.store == nil {
		return fmt.Errorf("config: no storage driver")
	}

	data, err := m.store.Load()
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return ErrNoConfig
	}

	if err := cbor.Unmarshal(data, v); err != nil {
		logger.Error("Config: Corrupt data found. Decode failed: %v", err)
		return fmt.Errorf("config decode failed: %w", err)
	}

	if validator, ok := v.(Validator); ok {
		if err := validator.Validate(); err != nil {
			logger.Error("Config: Validation failed: %v", err)
			return fmt.Errorf("config validation failed: %w", err)
		}
	}

	logger.Debug("Config loaded (%d bytes)", len(data))
	return nil
}

// Save persists configuration.
func (m *Manager) Save(v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return fmt.Errorf("config: no storage driver")
	}

	data, err := cbor.Marshal(v)
	if err != nil {
		return fmt.Errorf("config encode failed: %w", err)
	}

	if err := m.store.Save(data); err != nil {
		logger.Error("Config: Write failed: %v", err)
		return err
	}

	logger.Info("Config saved (%d bytes)", len(data))
	return nil
}

// Reset clears the configuration.
func (m *Manager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.store == nil {
		return nil
	}
	
	logger.Warn("Config: Performing factory reset")
	return m.store.Clear()
}