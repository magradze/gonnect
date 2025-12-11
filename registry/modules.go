// registry/modules.go
package registry

import (
	"fmt"
	"sync"

	"github.com/magradze/gonnect"
	"github.com/magradze/gonnect/pkg/logger"
)

var (
	modulesMu sync.Mutex
	// modules holds the list of registered components.
	// We use a Slice instead of a Map to preserve initialization order (Deterministic Startup).
	// The order is determined by the import order in main.go.
	modules []gonnect.Module
)

// RegisterModule adds a new module to the system lifecycle.
// This is typically called within the init() function of the module package.
// It panics if a module with the same name is already registered.
func RegisterModule(m gonnect.Module) {
	modulesMu.Lock()
	defer modulesMu.Unlock()

	name := m.Name()
	if name == "" {
		panic("gonnect: attempted to register module with empty name")
	}

	// Linear scan for duplicates (O(N)).
	// Since N (module count) is small in embedded systems (<50),
	// this is faster and lighter on RAM than a Map hash lookup.
	for _, existing := range modules {
		if existing.Name() == name {
			panic(fmt.Sprintf("gonnect: module '%s' is already registered", name))
		}
	}

	modules = append(modules, m)
	logger.Debug("Module registered: '%s'", name)
}

// GetModules returns the slice of registered modules.
// The Engine uses this list to boot components in the order they were imported.
func GetModules() []gonnect.Module {
	modulesMu.Lock()
	defer modulesMu.Unlock()
	return modules
}