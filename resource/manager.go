package resource

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect/pkg/logger"
)

// Manager handles the allocation and locking of hardware resources.
// It ensures that no two drivers attempt to use the same physical resource simultaneously.
type Manager struct {
	mu    sync.Mutex
	locks map[Type]map[ID]string // map[ResourceType][ResourceID]OwnerName
}

// globalManager is the singleton instance of the resource manager.
var globalManager = &Manager{
	locks: make(map[Type]map[ID]string),
}

// Lock attempts to claim a specific resource for a given owner.
// It returns an error if the resource is already locked by another component.
// This function is thread-safe.
func Lock(t Type, id ID, owner string) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	// Initialize map for this resource type if it doesn't exist
	if _, ok := globalManager.locks[t]; !ok {
		globalManager.locks[t] = make(map[ID]string)
	}

	// Check if the resource is already locked
	if currentOwner, exists := globalManager.locks[t][id]; exists {
		errMsg := fmt.Sprintf("resource conflict: %s/%d is already locked by '%s', requested by '%s'",
			t, id, currentOwner, owner)
		logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Lock the resource
	globalManager.locks[t][id] = owner
	logger.Debug("Resource locked: %s/%d by '%s'", t, id, owner)

	return nil
}

// Unlock releases a previously claimed resource.
// It returns an error if the resource was not locked or if the requester is not the owner.
// This function is thread-safe.
func Unlock(t Type, id ID, owner string) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	// Check if the resource type map exists
	if _, ok := globalManager.locks[t]; !ok {
		return fmt.Errorf("resource type %s has no active locks", t)
	}

	// Check if the specific resource is locked
	currentOwner, exists := globalManager.locks[t][id]
	if !exists {
		return fmt.Errorf("resource %s/%d is not locked", t, id)
	}

	// Verify ownership
	if currentOwner != owner {
		errMsg := fmt.Sprintf("security violation: '%s' attempted to unlock %s/%d owned by '%s'",
			owner, t, id, currentOwner)
		logger.Warn(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Unlock the resource
	delete(globalManager.locks[t], id)
	logger.Debug("Resource unlocked: %s/%d by '%s'", t, id, owner)

	return nil
}

// IsLocked checks if a specific resource is currently in use.
// Useful for diagnostics and debugging.
func IsLocked(t Type, id ID) bool {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if typeMap, ok := globalManager.locks[t]; ok {
		_, exists := typeMap[id]
		return exists
	}
	return false
}

// GetOwner returns the name of the component holding the lock, if any.
// Returns an empty string if the resource is free.
func GetOwner(t Type, id ID) string {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if typeMap, ok := globalManager.locks[t]; ok {
		if owner, exists := typeMap[id]; exists {
			return owner
		}
	}
	return ""
}