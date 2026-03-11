package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// DNSUpstreamPoolSpec defines the desired state of DNSUpstreamPool.
type DNSUpstreamPoolSpec struct {
	// Upstreams is the list of upstream DNS resolvers.
	Upstreams []Upstream `json:"upstreams"`
	// HealthCheck configures health checking for upstreams.
	HealthCheck HealthCheckConfig `json:"healthCheck,omitempty"`
	// LoadBalancing configures how queries are distributed across upstreams.
	LoadBalancing LoadBalancingConfig `json:"loadBalancing,omitempty"`
}

// Upstream defines a single upstream DNS resolver.
type Upstream struct {
	// Address is the IP address or hostname of the upstream resolver.
	Address string `json:"address"`
	// Port is the port number of the upstream resolver. Defaults to 53.
	Port int32 `json:"port,omitempty"`
}

// HealthCheckConfig configures upstream health checking.
type HealthCheckConfig struct {
	// Enabled enables periodic health checks. Defaults to true.
	Enabled bool `json:"enabled,omitempty"`
	// IntervalSeconds is the time between health checks. Defaults to 30.
	IntervalSeconds int32 `json:"intervalSeconds,omitempty"`
	// TimeoutSeconds is the health check timeout. Defaults to 5.
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// FailureThreshold is the number of consecutive failures before marking unhealthy. Defaults to 3.
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// LoadBalancingConfig configures upstream load balancing.
type LoadBalancingConfig struct {
	// Strategy is the load balancing algorithm: round-robin, first-available, or random.
	Strategy string `json:"strategy,omitempty"`
}

// DNSUpstreamPoolStatus defines the observed state of DNSUpstreamPool.
type DNSUpstreamPoolStatus struct {
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
