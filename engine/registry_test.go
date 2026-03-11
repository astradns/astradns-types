package engine

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type testEngine struct {
	configDir string
}

func (e *testEngine) Configure(context.Context, EngineConfig) (string, error) { return "", nil }
func (e *testEngine) Start(context.Context) error                             { return nil }
func (e *testEngine) Reload(context.Context) error                            { return nil }
func (e *testEngine) Stop(context.Context) error                              { return nil }
func (e *testEngine) HealthCheck(context.Context) (bool, error)               { return true, nil }
func (e *testEngine) Name() EngineType                                        { return EngineType("test") }

func resetRegistryForTest(t *testing.T) {
	t.Helper()

	registryMu.Lock()
	previous := make(map[EngineType]EngineFactory, len(registry))
	for key, value := range registry {
		previous[key] = value
	}
	registry = map[EngineType]EngineFactory{}
	registryMu.Unlock()

	t.Cleanup(func() {
		registryMu.Lock()
		registry = previous
		registryMu.Unlock()
	})
}

func TestRegisterAndNew(t *testing.T) {
	resetRegistryForTest(t)

	Register(EngineType("custom"), func(configDir string) Engine {
		return &testEngine{configDir: configDir}
	})

	created, err := New(EngineType("custom"), "/tmp/engine")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	test, ok := created.(*testEngine)
	if !ok {
		t.Fatalf("expected *testEngine, got %T", created)
	}
	if test.configDir != "/tmp/engine" {
		t.Fatalf("expected config dir %q, got %q", "/tmp/engine", test.configDir)
	}
}

func TestAvailableEnginesUsesDeterministicOrder(t *testing.T) {
	resetRegistryForTest(t)

	Register(EngineType("zeta"), func(configDir string) Engine { return &testEngine{configDir: configDir} })
	Register(EnginePowerDNS, func(configDir string) Engine { return &testEngine{configDir: configDir} })
	Register(EngineType("alpha"), func(configDir string) Engine { return &testEngine{configDir: configDir} })
	Register(EngineUnbound, func(configDir string) Engine { return &testEngine{configDir: configDir} })
	Register(EngineCoreDNS, func(configDir string) Engine { return &testEngine{configDir: configDir} })

	got := AvailableEngines()
	want := []EngineType{EngineUnbound, EngineCoreDNS, EnginePowerDNS, EngineType("alpha"), EngineType("zeta")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected engine order\nwant: %v\n got: %v", want, got)
	}
}

func TestRegisterConcurrentAccessIsSafe(t *testing.T) {
	resetRegistryForTest(t)

	const workers = 64
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(index int) {
			defer wg.Done()
			engineType := EngineType(fmt.Sprintf("engine-%d", index))
			Register(engineType, func(configDir string) Engine { return &testEngine{configDir: configDir} })
			_, _ = New(engineType, "/tmp")
			_ = AvailableEngines()
		}(i)
	}

	wg.Wait()

	if got := len(AvailableEngines()); got != workers {
		t.Fatalf("expected %d engines registered, got %d", workers, got)
	}
}
