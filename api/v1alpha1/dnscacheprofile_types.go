package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// DNSCacheProfileSpec defines the desired state of DNSCacheProfile.
type DNSCacheProfileSpec struct {
	// MaxEntries is the maximum number of cache entries. Defaults to 100000.
	// +kubebuilder:default=100000
	// +kubebuilder:validation:Minimum=1
	MaxEntries int32 `json:"maxEntries,omitempty"`
	// PositiveTtl configures TTL bounds for successful responses.
	PositiveTtl TtlConfig `json:"positiveTtl,omitempty"`
	// NegativeTtl configures TTL for negative responses (NXDOMAIN).
	NegativeTtl NegTtlConfig `json:"negativeTtl,omitempty"`
	// Prefetch configures cache prefetching.
	Prefetch PrefetchConfig `json:"prefetch,omitempty"`
}

// TtlConfig configures TTL bounds for cache entries.
// +kubebuilder:validation:XValidation:rule="self.maxSeconds == 0 || self.minSeconds == 0 || self.maxSeconds >= self.minSeconds",message="maxSeconds must be greater than or equal to minSeconds"
type TtlConfig struct {
	// MinSeconds is the minimum TTL in seconds. Defaults to 60.
	// +kubebuilder:default=60
	// +kubebuilder:validation:Minimum=1
	MinSeconds int32 `json:"minSeconds,omitempty"`
	// MaxSeconds is the maximum TTL in seconds. Defaults to 300.
	// +kubebuilder:default=300
	// +kubebuilder:validation:Minimum=1
	MaxSeconds int32 `json:"maxSeconds,omitempty"`
}

// NegTtlConfig configures negative response caching.
type NegTtlConfig struct {
	// Seconds is the TTL for negative cache entries. Defaults to 30.
	// +kubebuilder:default=30
	// +kubebuilder:validation:Minimum=1
	Seconds int32 `json:"seconds,omitempty"`
}

// PrefetchConfig configures cache prefetching behavior.
type PrefetchConfig struct {
	// Enabled enables cache prefetching.
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`
	// Threshold is the number of queries before a record is prefetched. Defaults to 10.
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=1
	Threshold int32 `json:"threshold,omitempty"`
}

// DNSCacheProfileStatus defines the observed state of DNSCacheProfile.
type DNSCacheProfileStatus struct {
	// ObservedGeneration is the most recent generation acted on by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions represent the latest available observations of the profile's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dcp
// +kubebuilder:printcolumn:name="Active",type=string,JSONPath=`.status.conditions[?(@.type=="Active")].status`
// +kubebuilder:printcolumn:name="MaxEntries",type=integer,JSONPath=`.spec.maxEntries`

// DNSCacheProfile configures DNS cache behavior.
type DNSCacheProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DNSCacheProfileSpec   `json:"spec,omitempty"`
	Status            DNSCacheProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DNSCacheProfileList contains a list of DNSCacheProfile.
type DNSCacheProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSCacheProfile `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSCacheProfile{}, &DNSCacheProfileList{})
}
