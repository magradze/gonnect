// resource/manager.go
package resource

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/logger"
)

// resourceKey acts as a composite unique identifier for the hash map.
// It avoids the overhead of nested maps (map[Type]map[ID]).
type resourceKey struct {
	Type Type
	ID   ID
}

// Manager handles the atomic allocation and locking of hardware resources.
type Manager struct {
	mu    sync.Mutex
	// locks stores the owner of each resource.
	// We use a flat map with a struct key to reduce heap allocations and GC scan time.
	locks map[resourceKey]string
}

// globalManager is the singleton instance.
// The map is lazily initialized to save memory if no resources are ever locked.
var globalManager = &Manager{}

// ensureInit initializes the map if it hasn't been created yet.
// This allows the binary to start with zero heap allocation for the manager.
func (m *Manager) ensureInit() {
	if m.locks == nil {
		m.locks = make(map[resourceKey]string)
	}
}

// Lock claims exclusive access to a hardware resource.
// It returns an error if the resource is already owned by another component.
func Lock(t Type, id ID, owner string) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	globalManager.ensureInit()

	key := resourceKey{Type: t, ID: id}

	// Fast path: Check for existence
	if currentOwner, exists := globalManager.locks[key]; exists {
		// Error construction is deferred until failure to avoid allocation on the happy path.
		errMsg := fmt.Sprintf("resource conflict: %s/%d owned by '%s', requested by '%s'",
			t, id, currentOwner, owner)
		logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Success
	globalManager.locks[key] = owner
	logger.Debug("Resource locked: %s/%d by '%s'", t, id, owner)

	return nil
}

// Unlock releases a resource.
// It enforces strict ownership validation to prevent unauthorized release.
func Unlock(t Type, id ID, owner string) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if globalManager.locks == nil {
		return fmt.Errorf("resource manager: no locks active")
	}

	key := resourceKey{Type: t, ID: id}

	currentOwner, exists := globalManager.locks[key]
	if !exists {
		return fmt.Errorf("resource unlock failed: %s/%d is not locked", t, id)
	}

	if currentOwner != owner {
		errMsg := fmt.Sprintf("security violation: '%s' tried to unlock %s/%d owned by '%s'",
			owner, t, id, currentOwner)
		logger.Warn(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Delete removes the key from the map.
	// Note: In Go maps, this does not shrink the memory footprint immediately,
	// but marks the slot as empty for reuse.
	delete(globalManager.locks, key)
	logger.Debug("Resource unlocked: %s/%d by '%s'", t, id, owner)

	return nil
}

// IsLocked checks if a resource is currently busy.
func IsLocked(t Type, id ID) bool {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if globalManager.locks == nil {
		return false
	}

	key := resourceKey{Type: t, ID: id}
	_, exists := globalManager.locks[key]
	return exists
}

// GetOwner returns the owner name of a resource or empty string.
func GetOwner(t Type, id ID) string {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if globalManager.locks == nil {
		return ""
	}

	key := resourceKey{Type: t, ID: id}
	return globalManager.locks[key]
}