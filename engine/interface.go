package engine

import (
	"context"
	"time"
)

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

// UpstreamTransport identifies how upstream DNS queries are transported.
type UpstreamTransport string

const (
	// UpstreamTransportDNS uses plain DNS over UDP/TCP.
	UpstreamTransportDNS UpstreamTransport = "dns"
	// UpstreamTransportDoT uses DNS-over-TLS.
	UpstreamTransportDoT UpstreamTransport = "dot"
	// UpstreamTransportDoH uses DNS-over-HTTPS.
	UpstreamTransportDoH UpstreamTransport = "doh"
)

// DNSSECMode identifies DNSSEC processing behavior.
type DNSSECMode string

const (
	// DNSSECModeOff disables DNSSEC processing.
	DNSSECModeOff DNSSECMode = "off"
	// DNSSECModeProcess enables DNSSEC-aware processing without strict validation.
	DNSSECModeProcess DNSSECMode = "process"
	// DNSSECModeValidate enables strict DNSSEC validation.
	DNSSECModeValidate DNSSECMode = "validate"
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

	// Capabilities describes the feature surface supported by the engine implementation.
	Capabilities() EngineCapabilities

	// HealthStatus returns a detailed health snapshot for the running engine.
	HealthStatus(ctx context.Context) (EngineHealthStatus, error)

	// HealthCheck returns true if the engine is responding to DNS queries.
	// Prefer HealthStatus when callers need latency and failure reason details.
	HealthCheck(ctx context.Context) (bool, error)

	// Name returns the engine identifier.
	Name() EngineType
}

// EngineCapabilities describes supported behavior for an engine implementation.
type EngineCapabilities struct {
	SupportsHotReload         bool
	SupportedTransports       []UpstreamTransport
	SupportedDNSSECModes      []DNSSECMode
	SupportsTLSServerName     bool
	SupportsWeightedUpstreams bool
	SupportsPriorityUpstreams bool
}

// EngineHealthStatus reports liveness plus diagnostic context.
type EngineHealthStatus struct {
	Healthy bool
	Latency time.Duration
	Reason  string
}

// EngineConfig is the engine-agnostic configuration derived from CRDs.
type EngineConfig struct {
	// Upstreams is the list of upstream resolvers to forward to.
	Upstreams []UpstreamConfig `json:"upstreams"`

	// Cache holds cache tuning parameters.
	Cache CacheConfig `json:"cache"`

	// ListenAddr is the address the engine should listen on (e.g., "127.0.0.1").
	ListenAddr string `json:"listenAddr"`

	// ListenPort is the port the engine should listen on (e.g., 5354).
	ListenPort int32 `json:"listenPort"`

	// WorkerThreads is the number of engine worker threads.
	WorkerThreads int32 `json:"workerThreads,omitempty"`

	// DNSSEC holds DNSSEC processing settings.
	DNSSEC DNSSECConfig `json:"dnssec,omitempty"`
}

// UpstreamConfig represents a single upstream resolver.
type UpstreamConfig struct {
	// Address is the IP or hostname of the upstream.
	Address string `json:"address"`
	// Port is the port of the upstream.
	Port int32 `json:"port"`
	// Transport selects upstream protocol transport.
	Transport UpstreamTransport `json:"transport,omitempty"`
	// TLSServerName overrides SNI/hostname verification for TLS-based transports.
	TLSServerName string `json:"tlsServerName,omitempty"`
	// Weight is the relative upstream weight.
	Weight int32 `json:"weight,omitempty"`
	// Preference is the priority hint for ordering upstreams.
	Preference int32 `json:"preference,omitempty"`
}

// DNSSECConfig holds DNSSEC processing settings.
type DNSSECConfig struct {
	// Mode controls resolver DNSSEC behavior.
	Mode DNSSECMode `json:"mode,omitempty"`
}

// CacheConfig holds cache tuning parameters.
type CacheConfig struct {
	// MaxEntries is the maximum number of cache entries.
	MaxEntries int32 `json:"maxEntries"`
	// PositiveTtlMin is the minimum TTL for positive responses in seconds.
	PositiveTtlMin int32 `json:"positiveTtlMin"`
	// PositiveTtlMax is the maximum TTL for positive responses in seconds.
	PositiveTtlMax int32 `json:"positiveTtlMax"`
	// NegativeTtl is the TTL for negative responses in seconds.
	NegativeTtl int32 `json:"negativeTtl"`
	// PrefetchEnabled enables cache prefetching.
	PrefetchEnabled bool `json:"prefetchEnabled"`
	// PrefetchThreshold is the number of lookups before prefetching.
	PrefetchThreshold int32 `json:"prefetchThreshold"`
}
