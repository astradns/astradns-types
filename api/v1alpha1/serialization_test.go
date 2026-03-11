package v1alpha1

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDNSUpstreamPoolJSONRoundTrip(t *testing.T) {
	original := DNSUpstreamPool{
		TypeMeta: metav1.TypeMeta{APIVersion: "dns.astradns.com/v1alpha1", Kind: "DNSUpstreamPool"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "dns-system",
		},
		Spec: DNSUpstreamPoolSpec{
			Upstreams: []Upstream{{Address: "1.1.1.1", Port: 53}},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal DNSUpstreamPool: %v", err)
	}

	var decoded DNSUpstreamPool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal DNSUpstreamPool: %v", err)
	}

	if decoded.Spec.Upstreams[0].Address != original.Spec.Upstreams[0].Address {
		t.Fatalf("expected upstream address %q, got %q", original.Spec.Upstreams[0].Address, decoded.Spec.Upstreams[0].Address)
	}
}

func TestDNSCacheProfileJSONRoundTrip(t *testing.T) {
	original := DNSCacheProfile{
		TypeMeta: metav1.TypeMeta{APIVersion: "dns.astradns.com/v1alpha1", Kind: "DNSCacheProfile"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "dns-system",
		},
		Spec: DNSCacheProfileSpec{
			MaxEntries: 100000,
			PositiveTtl: TtlConfig{
				MinSeconds: 60,
				MaxSeconds: 300,
			},
			NegativeTtl: NegTtlConfig{Seconds: 30},
			Prefetch:    PrefetchConfig{Enabled: true, Threshold: 10},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal DNSCacheProfile: %v", err)
	}

	var decoded DNSCacheProfile
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal DNSCacheProfile: %v", err)
	}

	if decoded.Spec.PositiveTtl.MinSeconds != original.Spec.PositiveTtl.MinSeconds {
		t.Fatalf("expected min ttl %d, got %d", original.Spec.PositiveTtl.MinSeconds, decoded.Spec.PositiveTtl.MinSeconds)
	}
	if decoded.Spec.Prefetch.Enabled != original.Spec.Prefetch.Enabled {
		t.Fatalf("expected prefetch enabled %t, got %t", original.Spec.Prefetch.Enabled, decoded.Spec.Prefetch.Enabled)
	}
}

func TestExternalDNSPolicyJSONRoundTrip(t *testing.T) {
	original := ExternalDNSPolicy{
		TypeMeta: metav1.TypeMeta{APIVersion: "dns.astradns.com/v1alpha1", Kind: "ExternalDNSPolicy"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy-a",
			Namespace: "tenant-a",
		},
		Spec: ExternalDNSPolicySpec{
			Selector: PolicySelector{Namespaces: []string{"tenant-a"}},
			UpstreamPoolRef: ResourceRef{
				Name: "default",
			},
			CacheProfileRef: ResourceRef{Name: "fast"},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal ExternalDNSPolicy: %v", err)
	}

	var decoded ExternalDNSPolicy
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal ExternalDNSPolicy: %v", err)
	}

	if decoded.Spec.UpstreamPoolRef.Name != original.Spec.UpstreamPoolRef.Name {
		t.Fatalf("expected upstreamPoolRef.name %q, got %q", original.Spec.UpstreamPoolRef.Name, decoded.Spec.UpstreamPoolRef.Name)
	}
	if len(decoded.Spec.Selector.Namespaces) != 1 || decoded.Spec.Selector.Namespaces[0] != "tenant-a" {
		t.Fatalf("expected namespace selector to include tenant-a, got %#v", decoded.Spec.Selector.Namespaces)
	}
}
