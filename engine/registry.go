package engine

import (
	"fmt"
	"sort"
	"sync"
)

// EngineFactory is a function that creates a new Engine instance.
type EngineFactory func(configDir string) Engine

// registry holds all registered engine factories.
var (
	registryMu sync.RWMutex
	registry   = map[EngineType]EngineFactory{}
)

// Register adds an engine factory to the registry.
// Called by each engine package's init() function.
func Register(engineType EngineType, factory EngineFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	registry[engineType] = factory
}

// New creates a new Engine instance for the given type.
func New(engineType EngineType, configDir string) (Engine, error) {
	registryMu.RLock()
	factory, ok := registry[engineType]
	registryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unknown engine type: %s (available: %v)", engineType, AvailableEngines())
	}
	return factory(configDir), nil
}

// AvailableEngines returns all registered engine types.
func AvailableEngines() []EngineType {
	registryMu.RLock()
	registered := make(map[EngineType]struct{}, len(registry))
	for t := range registry {
		registered[t] = struct{}{}
	}
	registryMu.RUnlock()

	ordered := []EngineType{EngineUnbound, EngineCoreDNS, EnginePowerDNS}
	types := make([]EngineType, 0, len(registered))

	for _, t := range ordered {
		if _, ok := registered[t]; !ok {
			continue
		}
		types = append(types, t)
		delete(registered, t)
	}

	unknown := make([]EngineType, 0, len(registered))
	for t := range registered {
		unknown = append(unknown, t)
	}
	sort.Slice(unknown, func(i, j int) bool {
		return string(unknown[i]) < string(unknown[j])
	})
	types = append(types, unknown...)

	return types
}
