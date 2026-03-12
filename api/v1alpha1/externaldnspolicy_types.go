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
	// SplitHorizon optionally defines source-aware DNS views with zone-specific routing.
	// This models the API contract for split-horizon policy; enforcement is phased separately.
	SplitHorizon *SplitHorizonPolicy `json:"splitHorizon,omitempty"`
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
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// SplitHorizonPolicy describes source-aware DNS routing views.
type SplitHorizonPolicy struct {
	// Views is the ordered list of matching views. First match wins.
	// +kubebuilder:validation:MinItems=1
	Views []SplitHorizonView `json:"views"`
}

// SplitHorizonView describes a source-specific view.
type SplitHorizonView struct {
	// Name is the stable view identifier.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name"`
	// SourceCIDRs is the list of source CIDRs matched by this view.
	// +kubebuilder:validation:MinItems=1
	SourceCIDRs []string `json:"sourceCIDRs"`
	// Zones defines zone-to-upstream routing within this view.
	// +kubebuilder:validation:MinItems=1
	Zones []SplitHorizonZoneRule `json:"zones"`
}

// SplitHorizonZoneRule maps a DNS zone to upstream/cache references.
type SplitHorizonZoneRule struct {
	// Zone is the DNS suffix matched by this rule (for example "corp.example.com").
	// +kubebuilder:validation:MinLength=1
	Zone string `json:"zone"`
	// UpstreamPoolRef references the pool used for this zone.
	UpstreamPoolRef ResourceRef `json:"upstreamPoolRef"`
	// CacheProfileRef optionally references a cache profile for this zone.
	CacheProfileRef ResourceRef `json:"cacheProfileRef,omitempty"`
}

// ExternalDNSPolicyStatus defines the observed state of ExternalDNSPolicy.
type ExternalDNSPolicyStatus struct {
	// ObservedGeneration is the most recent generation acted on by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions represent the latest available observations of the policy's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// AppliedNodes is the number of nodes where this policy is applied.
	AppliedNodes int32 `json:"appliedNodes,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=edp
// +kubebuilder:printcolumn:name="Validated",type=string,JSONPath=`.status.conditions[?(@.type=="Validated")].status`
// +kubebuilder:printcolumn:name="AppliedNodes",type=integer,JSONPath=`.status.appliedNodes`

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
