package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// DNSCacheProfileSpec defines the desired state of DNSCacheProfile.
type DNSCacheProfileSpec struct {
	// MaxEntries is the maximum number of cache entries. Defaults to 100000.
	MaxEntries int32 `json:"maxEntries,omitempty"`
	// PositiveTtl configures TTL bounds for successful responses.
	PositiveTtl TtlConfig `json:"positiveTtl,omitempty"`
	// NegativeTtl configures TTL for negative responses (NXDOMAIN).
	NegativeTtl NegTtlConfig `json:"negativeTtl,omitempty"`
	// Prefetch configures cache prefetching.
	Prefetch PrefetchConfig `json:"prefetch,omitempty"`
}

// TtlConfig configures TTL bounds for cache entries.
type TtlConfig struct {
	// MinSeconds is the minimum TTL in seconds. Defaults to 60.
	MinSeconds int32 `json:"minSeconds,omitempty"`
	// MaxSeconds is the maximum TTL in seconds. Defaults to 300.
	MaxSeconds int32 `json:"maxSeconds,omitempty"`
}

// NegTtlConfig configures negative response caching.
type NegTtlConfig struct {
	// Seconds is the TTL for negative cache entries. Defaults to 30.
	Seconds int32 `json:"seconds,omitempty"`
}

// PrefetchConfig configures cache prefetching behavior.
type PrefetchConfig struct {
	// Enabled enables cache prefetching.
	Enabled bool `json:"enabled,omitempty"`
	// Threshold is the number of queries before a record is prefetched. Defaults to 10.
	Threshold int32 `json:"threshold,omitempty"`
}

// DNSCacheProfileStatus defines the observed state of DNSCacheProfile.
type DNSCacheProfileStatus struct {
	// Conditions represent the latest available observations of the profile's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
