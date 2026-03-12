package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// DNSUpstreamPoolSpec defines the desired state of DNSUpstreamPool.
type DNSUpstreamPoolSpec struct {
	// Upstreams is the list of upstream DNS resolvers.
	// +kubebuilder:validation:MinItems=1
	Upstreams []Upstream `json:"upstreams"`
	// HealthCheck configures health checking for upstreams.
	HealthCheck HealthCheckConfig `json:"healthCheck,omitempty"`
	// LoadBalancing configures how queries are distributed across upstreams.
	LoadBalancing LoadBalancingConfig `json:"loadBalancing,omitempty"`
	// Runtime configures engine runtime tuning.
	Runtime RuntimeConfig `json:"runtime,omitempty"`
	// DNSSEC configures resolver DNSSEC processing mode.
	DNSSEC DNSSECConfig `json:"dnssec,omitempty"`
}

// Upstream defines a single upstream DNS resolver.
type Upstream struct {
	// Address is the IP address or hostname of the upstream resolver.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern="^[A-Za-z0-9]([A-Za-z0-9.:-]*[A-Za-z0-9])?$"
	Address string `json:"address"`
	// Port is the port number of the upstream resolver.
	// When omitted, defaults are transport-specific (53 for dns, 853 for dot, 443 for doh).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port,omitempty"`
	// Transport selects upstream protocol transport: dns, dot, or doh.
	// +kubebuilder:validation:Enum=dns;dot;doh
	// +kubebuilder:default=dns
	Transport string `json:"transport,omitempty"`
	// TLSServerName overrides SNI/hostname verification for TLS-based transports.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern="^$|^[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?$"
	TLSServerName string `json:"tlsServerName,omitempty"`
	// Weight is the relative upstream weight.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	Weight int32 `json:"weight,omitempty"`
	// Preference is the upstream priority hint. Lower values are preferred.
	// +kubebuilder:default=100
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	Preference int32 `json:"preference,omitempty"`
}

// HealthCheckConfig configures upstream health checking.
type HealthCheckConfig struct {
	// Enabled enables periodic health checks. Defaults to true.
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
	// IntervalSeconds is the time between health checks. Defaults to 30.
	// +kubebuilder:default=30
	// +kubebuilder:validation:Minimum=1
	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`
	// TimeoutSeconds is the health check timeout. Defaults to 5.
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=1
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// FailureThreshold is the number of consecutive failures before marking unhealthy. Defaults to 3.
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=1
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// LoadBalancingConfig configures upstream load balancing.
type LoadBalancingConfig struct {
	// Strategy is the load balancing algorithm: round-robin, first-available, or random.
	// +kubebuilder:validation:Enum=round-robin;first-available;random
	// +kubebuilder:default=round-robin
	Strategy string `json:"strategy,omitempty"`
}

// RuntimeConfig configures engine runtime tuning.
type RuntimeConfig struct {
	// WorkerThreads is the number of worker threads to configure on supported engines.
	// When omitted, the runtime selects a CPU-based default.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=256
	WorkerThreads int32 `json:"workerThreads,omitempty"`
}

// DNSSECConfig configures resolver DNSSEC processing.
type DNSSECConfig struct {
	// Mode controls DNSSEC behavior for supported engines.
	// +kubebuilder:validation:Enum=off;process;validate
	// +kubebuilder:default=off
	Mode string `json:"mode,omitempty"`
}

// DNSUpstreamPoolStatus defines the observed state of DNSUpstreamPool.
type DNSUpstreamPoolStatus struct {
	// ObservedGeneration is the most recent generation acted on by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions represent the latest available observations of the pool's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// UpstreamStatuses contains the health status of each upstream.
	UpstreamStatuses []UpstreamStatus `json:"upstreamStatuses,omitempty"`
}

// UpstreamStatus represents the observed state of a single upstream.
type UpstreamStatus struct {
	// Address is the upstream address.
	Address string `json:"address"`
	// Healthy indicates whether the upstream is responding to health checks.
	Healthy bool `json:"healthy"`
	// LatencyMs is the last measured latency in milliseconds.
	LatencyMs int64 `json:"latencyMs,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dnsup
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Strategy",type=string,JSONPath=`.spec.loadBalancing.strategy`

// DNSUpstreamPool is a pool of upstream DNS resolvers.
type DNSUpstreamPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DNSUpstreamPoolSpec   `json:"spec,omitempty"`
	Status            DNSUpstreamPoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DNSUpstreamPoolList contains a list of DNSUpstreamPool.
type DNSUpstreamPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSUpstreamPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSUpstreamPool{}, &DNSUpstreamPoolList{})
}
