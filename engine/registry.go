package engine

import "fmt"

// EngineFactory is a function that creates a new Engine instance.
type EngineFactory func(configDir string) Engine

// registry holds all registered engine factories.
var registry = map[EngineType]EngineFactory{}

// Register adds an engine factory to the registry.
// Called by each engine package's init() function.
func Register(engineType EngineType, factory EngineFactory) {
	registry[engineType] = factory
}

// New creates a new Engine instance for the given type.
func New(engineType EngineType, configDir string) (Engine, error) {
	factory, ok := registry[engineType]
	if !ok {
		return nil, fmt.Errorf("unknown engine type: %s (available: %v)", engineType, AvailableEngines())
	}
	return factory(configDir), nil
}

// AvailableEngines returns all registered engine types.
func AvailableEngines() []EngineType {
	types := make([]EngineType, 0, len(registry))
	for t := range registry {
		types = append(types, t)
	}
	return types
}
