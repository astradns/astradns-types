package engine

import "context"

// EngineType identifies which DNS engine to use.
type EngineType string

const (
	// EngineUnbound is the Unbound DNS resolver.
	EngineUnbound EngineType = "unbound"
	// EngineCoreDNS is the CoreDNS resolver.
	EngineCoreDNS EngineType = "coredns"
	// EnginePowerDNS is the PowerDNS Recursor.
	EnginePowerDNS EngineType = "powerdns"
)

// Engine manages the lifecycle of a DNS resolver engine.
type Engine interface {
	// Configure generates engine-specific config from the abstract EngineConfig.
	// Returns the path to the generated config file.
	Configure(ctx context.Context, config EngineConfig) (string, error)

	// Start launches the engine subprocess.
	Start(ctx context.Context) error

	// Reload triggers a graceful config reload without dropping queries.
	Reload(ctx context.Context) error

	// Stop gracefully shuts down the engine.
	Stop(ctx context.Context) error

	// HealthCheck returns true if the engine is responding to DNS queries.
	HealthCheck(ctx context.Context) (bool, error)

	// Name returns the engine identifier.
	Name() EngineType
}

// EngineConfig is the engine-agnostic configuration derived from CRDs.
type EngineConfig struct {
	// Upstreams is the list of upstream resolvers to forward to.
	Upstreams []UpstreamConfig

	// Cache holds cache tuning parameters.
	Cache CacheConfig

	// ListenAddr is the address the engine should listen on (e.g., "127.0.0.1").
	ListenAddr string

	// ListenPort is the port the engine should listen on (e.g., 5354).
	ListenPort int
}

// UpstreamConfig represents a single upstream resolver.
type UpstreamConfig struct {
	// Address is the IP or hostname of the upstream.
	Address string
	// Port is the port of the upstream.
	Port int
}

// CacheConfig holds cache tuning parameters.
type CacheConfig struct {
	// MaxEntries is the maximum number of cache entries.
	MaxEntries int
	// PositiveTtlMin is the minimum TTL for positive responses in seconds.
	PositiveTtlMin int
	// PositiveTtlMax is the maximum TTL for positive responses in seconds.
	PositiveTtlMax int
	// NegativeTtl is the TTL for negative responses in seconds.
	NegativeTtl int
	// PrefetchEnabled enables cache prefetching.
	PrefetchEnabled bool
	// PrefetchThreshold is the number of lookups before prefetching.
	PrefetchThreshold int
}
