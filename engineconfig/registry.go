package engineconfig

import (
	"fmt"

	"github.com/astradns/astradns-types/engine"
)

// RendererFactory creates a ConfigRenderer.
type RendererFactory func() ConfigRenderer

var rendererRegistry = map[engine.EngineType]RendererFactory{}

// RegisterRenderer adds a config renderer factory.
func RegisterRenderer(engineType engine.EngineType, factory RendererFactory) {
	rendererRegistry[engineType] = factory
}

// NewRenderer creates a ConfigRenderer for the given engine type.
func NewRenderer(engineType engine.EngineType) (ConfigRenderer, error) {
	factory, ok := rendererRegistry[engineType]
	if !ok {
		return nil, fmt.Errorf("no config renderer for engine: %s", engineType)
	}
	return factory(), nil
}
