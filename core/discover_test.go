package core

import (
	"context"
	"errors"
	"net"
	"testing"
)

func TestDiscoverFallsBackToWellKnown(t *testing.T) {
	orig := lookupSRV
	t.Cleanup(func() { lookupSRV = orig })

	lookupSRV = func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		return "", nil, errors.New("no such host")
	}

	url, err := Discover(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("Discover: unexpected error: %v", err)
	}
	want := "https://example.com/.well-known/jmap"
	if url != want {
		t.Fatalf("Discover = %q, want %q", url, want)
	}
}

func TestDiscoverFallsBackWhenSRVEmpty(t *testing.T) {
	orig := lookupSRV
	t.Cleanup(func() { lookupSRV = orig })

	lookupSRV = func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		return "", []*net.SRV{}, nil
	}

	url, err := Discover(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("Discover: unexpected error: %v", err)
	}
	want := "https://example.com/.well-known/jmap"
	if url != want {
		t.Fatalf("Discover = %q, want %q", url, want)
	}
}

func TestDiscoverUsesSRV(t *testing.T) {
	orig := lookupSRV
	t.Cleanup(func() { lookupSRV = orig })

	lookupSRV = func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		if service != "jmap" || proto != "tcp" || name != "example.com" {
			t.Errorf("lookupSRV(%q, %q, %q)", service, proto, name)
		}
		return "", []*net.SRV{{Target: "jmap.example.com.", Port: 443}}, nil
	}

	url, err := Discover(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("Discover: unexpected error: %v", err)
	}
	want := "https://jmap.example.com:443/.well-known/jmap"
	if url != want {
		t.Fatalf("Discover = %q, want %q", url, want)
	}
}

func TestDiscoverRejectsNonHostnameDomain(t *testing.T) {
	orig := lookupSRV
	t.Cleanup(func() { lookupSRV = orig })
	lookupSRV = func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		t.Fatal("lookupSRV must not run for an invalid domain")
		return "", nil, nil
	}

	cases := []string{
		"evil.example/steal",
		"user:pass@evil.example",
		"https://evil.example",
		"evil.example?x=1",
		"evil.example#frag",
		"evil.example/path",
		"",
		"evil.example:443",
	}
	for _, domain := range cases {
		_, err := Discover(context.Background(), domain)
		if err == nil {
			t.Fatalf("Discover(%q) succeeded, want error", domain)
		}
	}
}

func TestDiscoverContextCanceled(t *testing.T) {
	orig := lookupSRV
	t.Cleanup(func() { lookupSRV = orig })

	lookupSRV = func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
		t.Fatal("lookupSRV should not be called when context is already canceled")
		return "", nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Discover(ctx, "example.com")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Discover error = %v, want context.Canceled", err)
	}
}
