package engineconfig

import (
	v1alpha1 "github.com/astradns/astradns-types/api/v1alpha1"
	"github.com/astradns/astradns-types/engine"
)

// ConfigGenerator converts CRD objects to an engine-agnostic EngineConfig.
type ConfigGenerator interface {
	Generate(pool *v1alpha1.DNSUpstreamPool, profile *v1alpha1.DNSCacheProfile) (*engine.EngineConfig, error)
}

// ConfigRenderer converts an EngineConfig to an engine-specific config string.
type ConfigRenderer interface {
	// Render produces the engine-specific config string from an EngineConfig.
	Render(config *engine.EngineConfig) (string, error)
	// EngineType returns which engine this renderer is for.
	EngineType() engine.EngineType
	// ConfigFileName returns the config file name (e.g., "unbound.conf", "Corefile", "recursor.conf").
	ConfigFileName() string
}
