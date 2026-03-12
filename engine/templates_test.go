package engine

import (
	"math"
	"strings"
	"testing"
	"text/template"
)

func TestValidateTemplateConfigRejectsInjectionCharacters(t *testing.T) {
	config := EngineConfig{
		ListenAddr: "127.0.0.1",
		Upstreams: []UpstreamConfig{
			{Address: "1.1.1.1\nforward-addr: 6.6.6.6"},
		},
	}

	if err := ValidateTemplateConfig(config); err == nil {
		t.Fatal("expected validation error for newline injection")
	}
}

func TestValidateTemplateConfigAcceptsValidAddresses(t *testing.T) {
	config := EngineConfig{
		ListenAddr: "::1",
		Upstreams: []UpstreamConfig{
			{Address: "resolver1.example.com"},
			{Address: "8.8.8.8"},
		},
	}

	if err := ValidateTemplateConfig(config); err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
}

func TestNewTemplateDataClampsLargeCacheSizes(t *testing.T) {
	data := NewTemplateData(EngineConfig{
		ListenAddr: "127.0.0.1",
		Upstreams:  []UpstreamConfig{{Address: "1.1.1.1"}},
		Cache: CacheConfig{
			MaxEntries: math.MaxInt32,
		},
	})

	if data.MsgCacheSize != "256m" {
		t.Fatalf("expected msg cache size to clamp at 256m, got %q", data.MsgCacheSize)
	}
	if data.RrsetCacheSize != "512m" {
		t.Fatalf("expected rrset cache size to clamp at 512m, got %q", data.RrsetCacheSize)
	}
}

func TestFormatBytesRoundsUp(t *testing.T) {
	if got := formatBytes(1530920); got != "2m" {
		t.Fatalf("expected rounded-up value 2m, got %q", got)
	}
	if got := formatBytes(1537); got != "2k" {
		t.Fatalf("expected rounded-up value 2k, got %q", got)
	}
}

func TestNewTemplateDataNormalizesAddressesAndPorts(t *testing.T) {
	data := NewTemplateData(EngineConfig{
		ListenAddr: " 127.0.0.1 ",
		Upstreams: []UpstreamConfig{
			{Address: " 1.1.1.1 ", Port: 0},
			{Address: " dns.google ", Port: 5353},
		},
	})

	if data.ListenAddr != "127.0.0.1" {
		t.Fatalf("expected listen address to be trimmed, got %q", data.ListenAddr)
	}
	if got := data.Upstreams[0].Address; got != "1.1.1.1" {
		t.Fatalf("expected upstream address to be trimmed, got %q", got)
	}
	if got := data.ForwardAddresses; got != "1.1.1.1:53;dns.google:5353" {
		t.Fatalf("unexpected forward addresses: %q", got)
	}
}

func TestRecursorTemplateIncludesMinimumTTLOverride(t *testing.T) {
	tmpl, err := template.New("recursor.conf").Parse(RecursorConfTemplate)
	if err != nil {
		t.Fatalf("failed to parse recursor template: %v", err)
	}

	data := NewTemplateData(EngineConfig{
		Upstreams: []UpstreamConfig{{Address: "1.1.1.1", Port: 53}},
		Cache: CacheConfig{
			PositiveTtlMin: 60,
			PositiveTtlMax: 300,
			NegativeTtl:    30,
		},
		ListenAddr: "127.0.0.1",
		ListenPort: 5354,
	})

	var rendered strings.Builder
	if err := tmpl.Execute(&rendered, data); err != nil {
		t.Fatalf("failed to execute recursor template: %v", err)
	}

	if !strings.Contains(rendered.String(), "minimum-ttl-override=60") {
		t.Fatalf("expected minimum-ttl-override in rendered config:\n%s", rendered.String())
	}
}

func TestNewTemplateDataBuildsTransportSpecificCoreDNSUpstreams(t *testing.T) {
	data := NewTemplateData(EngineConfig{
		Upstreams: []UpstreamConfig{
			{Address: "1.1.1.1", Transport: UpstreamTransportDNS},
			{Address: "dns.quad9.net", Transport: UpstreamTransportDoT},
			{Address: "dns.google", Transport: UpstreamTransportDoH},
		},
	})

	if len(data.CoreDNSUpstreams) != 3 {
		t.Fatalf("expected 3 coredns upstream targets, got %d", len(data.CoreDNSUpstreams))
	}
	if got := data.CoreDNSUpstreams[0]; got != "1.1.1.1:53" {
		t.Fatalf("unexpected plain DNS upstream target: %q", got)
	}
	if got := data.CoreDNSUpstreams[1]; got != "tls://dns.quad9.net:853" {
		t.Fatalf("unexpected DoT upstream target: %q", got)
	}
	if got := data.CoreDNSUpstreams[2]; got != "https://dns.google:443" {
		t.Fatalf("unexpected DoH upstream target: %q", got)
	}
}

func TestNewTemplateDataSetsUnboundForwardTLSUpstream(t *testing.T) {
	t.Run("all dot upstreams enable unbound tls forwarding", func(t *testing.T) {
		data := NewTemplateData(EngineConfig{
			Upstreams: []UpstreamConfig{
				{Address: "dns.quad9.net", Transport: UpstreamTransportDoT},
				{Address: "dns.google", Transport: UpstreamTransportDoT},
			},
		})

		if !data.UnboundForwardTLSUpstream {
			t.Fatal("expected unbound forward-tls-upstream to be enabled")
		}
	})

	t.Run("mixed transports disable unbound tls forwarding", func(t *testing.T) {
		data := NewTemplateData(EngineConfig{
			Upstreams: []UpstreamConfig{
				{Address: "1.1.1.1", Transport: UpstreamTransportDNS},
				{Address: "dns.quad9.net", Transport: UpstreamTransportDoT},
			},
		})

		if data.UnboundForwardTLSUpstream {
			t.Fatal("expected unbound forward-tls-upstream to be disabled for mixed transports")
		}
	})
}

func TestValidateTemplateConfigRejectsInvalidTransport(t *testing.T) {
	err := ValidateTemplateConfig(EngineConfig{
		ListenAddr: "127.0.0.1",
		Upstreams:  []UpstreamConfig{{Address: "1.1.1.1", Transport: "bogus"}},
	})
	if err == nil {
		t.Fatal("expected invalid transport error")
	}
}

func TestValidateTemplateConfigRejectsTLSServerNameOnPlainDNS(t *testing.T) {
	err := ValidateTemplateConfig(EngineConfig{
		ListenAddr: "127.0.0.1",
		Upstreams:  []UpstreamConfig{{Address: "1.1.1.1", Transport: UpstreamTransportDNS, TLSServerName: "dns.example"}},
	})
	if err == nil {
		t.Fatal("expected tlsServerName validation error on plain DNS transport")
	}
}
