package engineconfig

import (
	"fmt"
	"sync"
	"testing"

	"github.com/astradns/astradns-types/engine"
)

type testRenderer struct {
	engineType engine.EngineType
}

func (r *testRenderer) Render(config *engine.EngineConfig) (string, error) {
	if config == nil {
		return "", fmt.Errorf("config is required")
	}
	return "ok", nil
}

func (r *testRenderer) EngineType() engine.EngineType { return r.engineType }
func (r *testRenderer) ConfigFileName() string        { return "test.conf" }

func resetRendererRegistryForTest(t *testing.T) {
	t.Helper()

	rendererRegistryMu.Lock()
	previous := make(map[engine.EngineType]RendererFactory, len(rendererRegistry))
	for key, value := range rendererRegistry {
		previous[key] = value
	}
	rendererRegistry = map[engine.EngineType]RendererFactory{}
	rendererRegistryMu.Unlock()

	t.Cleanup(func() {
		rendererRegistryMu.Lock()
		rendererRegistry = previous
		rendererRegistryMu.Unlock()
	})
}

func TestRegisterRendererAndNewRenderer(t *testing.T) {
	resetRendererRegistryForTest(t)

	engineType := engine.EngineType("custom")
	RegisterRenderer(engineType, func() ConfigRenderer {
		return &testRenderer{engineType: engineType}
	})

	renderer, err := NewRenderer(engineType)
	if err != nil {
		t.Fatalf("NewRenderer returned error: %v", err)
	}
	if renderer.EngineType() != engineType {
		t.Fatalf("expected engine type %q, got %q", engineType, renderer.EngineType())
	}
}

func TestRegisterRendererConcurrentAccessIsSafe(t *testing.T) {
	resetRendererRegistryForTest(t)

	const workers = 64
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(index int) {
			defer wg.Done()
			engineType := engine.EngineType(fmt.Sprintf("renderer-%d", index))
			RegisterRenderer(engineType, func() ConfigRenderer {
				return &testRenderer{engineType: engineType}
			})
			_, _ = NewRenderer(engineType)
		}(i)
	}

	wg.Wait()

	rendererRegistryMu.RLock()
	defer rendererRegistryMu.RUnlock()
	if len(rendererRegistry) != workers {
		t.Fatalf("expected %d registered renderers, got %d", workers, len(rendererRegistry))
	}
}
