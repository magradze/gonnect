package registry

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect"
	"github.com/magradze/gonnect/pkg/logger"
)

var (
	modulesMu sync.Mutex
	// modules holds the set of registered components.
	// We use a map to enforce unique naming immediately upon registration.
	modules = make(map[string]gonnect.Module)
)

// RegisterModule adds a new module to the system lifecycle.
// This is typically called within the init() function of the module package.
// It panics if a module with the same name is already registered, as this indicates
// a critical configuration error during build time.
func RegisterModule(m gonnect.Module) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	name := m.Name()
	if name == "" {
		panic("gonnect: attempted to register module with empty name")
	}

	if _, exists := modules[name]; exists {
		panic(fmt.Sprintf("gonnect: module '%s' is already registered", name))
	}

	modules[name] = m
	logger.Debug("Module registered: '%s'", name)
}

// GetModules returns a slice of all registered modules.
// The Engine uses this list to iterate through components during boot and shutdown sequences.
// Note: The order of iteration is random due to the underlying map implementation.
// If initialization order matters, modules should handle dependencies explicitly via the Service Locator.
func GetModules() []gonnect.Module {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	list := make([]gonnect.Module, 0, len(modules))
	for _, m := range modules {
		list = append(list, m)
	}
	return list
}