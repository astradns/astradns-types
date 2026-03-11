package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ExternalDNSPolicySpec defines the desired state of ExternalDNSPolicy.
type ExternalDNSPolicySpec struct {
	// Selector determines which namespaces this policy applies to.
	Selector PolicySelector `json:"selector"`
	// UpstreamPoolRef references the DNSUpstreamPool to use.
	UpstreamPoolRef ResourceRef `json:"upstreamPoolRef"`
	// CacheProfileRef optionally references a DNSCacheProfile.
	CacheProfileRef ResourceRef `json:"cacheProfileRef,omitempty"`
}

// PolicySelector selects target namespaces for the policy.
type PolicySelector struct {
	// Namespaces is the list of namespace names this policy targets.
	// +kubebuilder:validation:MinItems=1
	Namespaces []string `json:"namespaces"`
}

// ResourceRef is a reference to another resource by name.
type ResourceRef struct {
	// Name is the name of the referenced resource.
	Name string `json:"name"`
}

// ExternalDNSPolicyStatus defines the observed state of ExternalDNSPolicy.
type ExternalDNSPolicyStatus struct {
	// Conditions represent the latest available observations of the policy's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// AppliedNodes is the number of nodes where this policy is applied.
	AppliedNodes int32 `json:"appliedNodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ExternalDNSPolicy defines per-namespace DNS resolution policy.
type ExternalDNSPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ExternalDNSPolicySpec   `json:"spec,omitempty"`
	Status            ExternalDNSPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExternalDNSPolicyList contains a list of ExternalDNSPolicy.
type ExternalDNSPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalDNSPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalDNSPolicy{}, &ExternalDNSPolicyList{})
}
