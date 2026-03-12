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
			Upstreams: []Upstream{{
				Address:       "1.1.1.1",
				Port:          853,
				Transport:     "dot",
				TLSServerName: "dns.quad9.net",
				Weight:        5,
				Preference:    10,
			}},
			Runtime: RuntimeConfig{WorkerThreads: 4},
			DNSSEC:  DNSSECConfig{Mode: "validate"},
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
	if decoded.Spec.Upstreams[0].Transport != "dot" {
		t.Fatalf("expected upstream transport dot, got %q", decoded.Spec.Upstreams[0].Transport)
	}
	if decoded.Spec.Runtime.WorkerThreads != 4 {
		t.Fatalf("expected runtime workerThreads 4, got %d", decoded.Spec.Runtime.WorkerThreads)
	}
	if decoded.Spec.DNSSEC.Mode != "validate" {
		t.Fatalf("expected dnssec mode validate, got %q", decoded.Spec.DNSSEC.Mode)
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
			SplitHorizon: &SplitHorizonPolicy{
				Views: []SplitHorizonView{{
					Name:        "corp-view",
					SourceCIDRs: []string{"10.0.0.0/8"},
					Zones: []SplitHorizonZoneRule{{
						Zone:            "corp.example.com",
						UpstreamPoolRef: ResourceRef{Name: "corp-upstream"},
					}},
				}},
			},
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
	if decoded.Spec.SplitHorizon == nil || len(decoded.Spec.SplitHorizon.Views) != 1 {
		t.Fatalf("expected split-horizon view to round-trip, got %#v", decoded.Spec.SplitHorizon)
	}
}
