package engineconfig

import (
	"fmt"
	"sync"

	"github.com/astradns/astradns-types/engine"
)

// RendererFactory creates a ConfigRenderer.
type RendererFactory func() ConfigRenderer

var (
	rendererRegistryMu sync.RWMutex
	rendererRegistry   = map[engine.EngineType]RendererFactory{}
)

// RegisterRenderer adds a config renderer factory.
func RegisterRenderer(engineType engine.EngineType, factory RendererFactory) {
	rendererRegistryMu.Lock()
	defer rendererRegistryMu.Unlock()

	rendererRegistry[engineType] = factory
}

// NewRenderer creates a ConfigRenderer for the given engine type.
func NewRenderer(engineType engine.EngineType) (ConfigRenderer, error) {
	rendererRegistryMu.RLock()
	factory, ok := rendererRegistry[engineType]
	rendererRegistryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no config renderer for engine: %s", engineType)
	}
	return factory(), nil
}
